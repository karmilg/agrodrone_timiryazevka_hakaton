const state = {
	selectedFieldId: null,
	field: null,
	problems: [],
	drones: [],
};

const $ = (id) => document.getElementById(id);

async function api(method, url, body) {
	const opts = { method, headers: { "Content-Type": "application/json" } };
	if (body !== undefined) opts.body = JSON.stringify(body);
	const res = await fetch(url, opts);
	const data = await res.json().catch(() => ({}));
	if (!res.ok) throw new Error(data.error || `${method} ${url} -> ${res.status}`);
	return data;
}

function fmt(num, d = 2) {
	return (Number(num) || 0).toFixed(d);
}

async function refresh() {
	try {
		const s = await api("GET", "/api/state");
		$("tick").textContent = s.tick;
		$("now").textContent = new Date(s.now).toLocaleString("ru-RU");

		state.drones = s.drones || [];
		renderFields(s.fields || []);
		renderDrones(state.drones);
		renderTasks(s.tasks || []);

		if (state.selectedFieldId !== null) {
			await selectField(state.selectedFieldId, true);
		}
	} catch (e) {
		appendLog("⚠ " + e.message);
	}
}

function renderFields(fields) {
	const ul = $("fields-list");
	ul.innerHTML = "";
	if (fields.length === 0) {
		ul.innerHTML = `<li class="empty">нет полей</li>`;
		return;
	}
	fields.sort((a, b) => b.problem_index - a.problem_index);

	for (const f of fields) {
		const li = document.createElement("li");
		if (f.id === state.selectedFieldId) li.classList.add("active");

		const ndvi = fmt(f.avg_ndvi);
		const badgeClass = f.avg_ndvi > 0.6 ? "green" : f.avg_ndvi > 0.3 ? "orange" : "red";

		li.innerHTML = `
			<div class="row-main">
				<span class="row-title">${f.name}</span>
				<span class="row-sub">${f.width}×${f.height} клеток</span>
			</div>
			<span class="badge ${badgeClass}">NDVI ${ndvi}</span>
		`;
		li.onclick = () => selectField(f.id);
		ul.appendChild(li);
	}
}

async function selectField(id, silent = false) {
	state.selectedFieldId = id;
	try {
		const [field, problems] = await Promise.all([
			api("GET", `/api/fields/${id}`),
			api("GET", `/api/fields/${id}/problems`),
		]);
		state.field = field;
		state.problems = problems.cells || [];
		$("field-title").textContent = `${field.Name} · ${field.Width}×${field.Height}`;
		drawField();

		if (!silent) {
			document.querySelectorAll("#fields-list li").forEach((li) => li.classList.remove("active"));
		}
	} catch (e) {
		appendLog("⚠ " + e.message);
	}
}

function drawField() {
	const canvas = $("field-map");
	const ctx = canvas.getContext("2d");
	ctx.clearRect(0, 0, canvas.width, canvas.height);

	if (!state.field) return;

	const f = state.field;
	const size = Math.floor(Math.min(canvas.width, canvas.height) / Math.max(f.Width, f.Height));
	const offX = Math.floor((canvas.width - size * f.Width) / 2);
	const offY = Math.floor((canvas.height - size * f.Height) / 2);

	for (let y = 0; y < f.Height; y++) {
		for (let x = 0; x < f.Width; x++) {
			const c = f.Grid[y][x];
			ctx.fillStyle = cellColor(c);
			ctx.fillRect(offX + x * size + 1, offY + y * size + 1, size - 2, size - 2);
		}
	}

	ctx.strokeStyle = "#fb923c";
	ctx.lineWidth = 2;
	for (const [x, y] of state.problems) {
		ctx.strokeRect(
			offX + x * size + 1,
			offY + y * size + 1,
			size - 2,
			size - 2,
		);
	}

	for (const d of state.drones) {
		if (d.X < 0 || d.Y < 0 || d.X >= f.Width || d.Y >= f.Height) continue;
		const cx = offX + d.X * size + size / 2;
		const cy = offY + d.Y * size + size / 2;
		const r = Math.max(4, size / 3);

		const grad = ctx.createRadialGradient(cx, cy, 0, cx, cy, r * 2);
		grad.addColorStop(0, "rgba(96,165,250,0.5)");
		grad.addColorStop(1, "rgba(96,165,250,0)");
		ctx.fillStyle = grad;
		ctx.beginPath();
		ctx.arc(cx, cy, r * 2, 0, Math.PI * 2);
		ctx.fill();

		ctx.fillStyle = "#60a5fa";
		ctx.beginPath();
		ctx.arc(cx, cy, r, 0, Math.PI * 2);
		ctx.fill();

		ctx.strokeStyle = "#1e3a8a";
		ctx.lineWidth = 1.5;
		ctx.stroke();
	}
}

function cellColor(c) {
	if (!c.Scanned) return "#1f2937";
	const t = Math.max(0, Math.min(1, c.NDVI));
	const r = Math.round(248 * (1 - t) + 74 * t);
	const g = Math.round(113 * (1 - t) + 222 * t);
	const b = Math.round(113 * (1 - t) + 128 * t);
	return `rgb(${r},${g},${b})`;
}

function renderDrones(drones) {
	const ul = $("drones-list");
	ul.innerHTML = "";
	if (drones.length === 0) {
		ul.innerHTML = `<li class="empty">нет дронов</li>`;
		return;
	}
	for (const d of drones) {
		const battColor = d.Battery > 50 ? "green" : d.Battery > 20 ? "orange" : "red";
		const li = document.createElement("li");
		li.innerHTML = `
			<div class="row-main">
				<span class="row-title">${d.Name}</span>
				<span class="row-sub">${d.Status} · (${d.X}, ${d.Y})</span>
			</div>
			<span class="badge ${battColor}">${d.Battery}%</span>
		`;
		ul.appendChild(li);
	}
}

function renderTasks(tasks) {
	const ul = $("tasks-list");
	ul.innerHTML = "";
	if (tasks.length === 0) {
		ul.innerHTML = `<li class="empty">очередь пуста</li>`;
		return;
	}
	tasks.sort((a, b) => b.Priority - a.Priority);
	for (const t of tasks) {
		const li = document.createElement("li");
		li.innerHTML = `
			<div class="row-main">
				<span class="row-title">#${t.ID} · ${t.Type}</span>
				<span class="row-sub">поле ${t.FieldID} → (${t.TargetX}, ${t.TargetY})</span>
			</div>
			<span class="badge blue">P${t.Priority}</span>
		`;
		ul.appendChild(li);
	}
}

async function doTick() {
	try {
		const res = await api("POST", "/api/tick");
		const ts = new Date(res.now).toLocaleTimeString("ru-RU");
		for (const line of res.log || []) appendLog(line, ts);
		if (!res.log || res.log.length === 0) appendLog("тик прошёл без событий", ts);
		await refresh();
	} catch (e) {
		appendLog("⚠ " + e.message);
	}
}

function appendLog(text, ts) {
	const log = $("log");
	const empty = log.querySelector(".empty-log");
	if (empty) empty.remove();

	const row = document.createElement("div");
	row.className = "row";
	const time = ts || new Date().toLocaleTimeString("ru-RU");
	row.innerHTML = `<span class="ts">${time}</span>${text}`;
	log.appendChild(row);
	log.scrollTop = log.scrollHeight;
}

function clearLog() {
	$("log").innerHTML = `<div class="empty-log">лог пуст</div>`;
}

const modal = {
	current: null,
	open(title, fields, onOk) {
		$("modal-title").textContent = title;
		const body = $("modal-body");
		body.innerHTML = "";
		for (const f of fields) {
			const wrap = document.createElement("div");
			wrap.className = "field-group";
			wrap.innerHTML = `<label>${f.label}</label>`;
			if (f.type === "select") {
				const sel = document.createElement("select");
				sel.id = "f-" + f.name;
				for (const opt of f.options) {
					const o = document.createElement("option");
					o.value = opt.value;
					o.textContent = opt.label;
					sel.appendChild(o);
				}
				wrap.appendChild(sel);
			} else {
				const input = document.createElement("input");
				input.type = f.type || "text";
				input.id = "f-" + f.name;
				if (f.value !== undefined) input.value = f.value;
				wrap.appendChild(input);
			}
			body.appendChild(wrap);
		}
		modal.current = { fields, onOk };
		$("modal").classList.remove("hidden");
	},
	close() {
		$("modal").classList.add("hidden");
		modal.current = null;
	},
	collect() {
		const out = {};
		for (const f of modal.current.fields) {
			const el = $("f-" + f.name);
			out[f.name] = f.type === "number" ? Number(el.value) : el.value;
		}
		return out;
	},
};

function openCreateField() {
	modal.open("Новое поле", [
		{ name: "name", label: "Название", value: "Поле" },
		{ name: "width", label: "Ширина (клеток)", type: "number", value: 10 },
		{ name: "height", label: "Высота (клеток)", type: "number", value: 10 },
	], async (data) => {
		await api("POST", "/api/fields", data);
		appendLog(`создано поле «${data.name}» ${data.width}×${data.height}`);
		await refresh();
	});
}

function openCreateDrone() {
	modal.open("Новый дрон", [
		{ name: "name", label: "Имя", value: "Стрекоза" },
		{
			name: "kind", label: "Тип", type: "select",
			options: [
				{ value: "scout", label: "Разведчик" },
				{ value: "spray", label: "Опрыскиватель" },
			],
		},
	], async (data) => {
		await api("POST", "/api/drones", data);
		appendLog(`добавлен дрон «${data.name}» (${data.kind})`);
		await refresh();
	});
}

function openCreateTask() {
	modal.open("Новая задача", [
		{
			name: "type", label: "Тип", type: "select",
			options: [
				{ value: "scan", label: "Сканирование" },
				{ value: "spray", label: "Опрыскивание" },
			],
		},
		{ name: "field_id", label: "ID поля", type: "number", value: 1 },
		{ name: "x", label: "Координата X", type: "number", value: 0 },
		{ name: "y", label: "Координата Y", type: "number", value: 0 },
		{ name: "priority", label: "Приоритет", type: "number", value: 5 },
	], async (data) => {
		await api("POST", "/api/tasks", data);
		appendLog(`поставлена задача ${data.type} на (${data.x},${data.y})`);
		await refresh();
	});
}

$("btn-tick").addEventListener("click", doTick);
$("btn-add-field").addEventListener("click", openCreateField);
$("btn-add-drone").addEventListener("click", openCreateDrone);
$("btn-add-task").addEventListener("click", openCreateTask);
$("btn-clear-log").addEventListener("click", clearLog);
$("modal-cancel").addEventListener("click", () => modal.close());
$("modal-ok").addEventListener("click", async () => {
	try {
		const data = modal.collect();
		const fn = modal.current.onOk;
		modal.close();
		await fn(data);
	} catch (e) {
		appendLog("⚠ " + e.message);
	}
});

$("modal").addEventListener("click", (e) => {
	if (e.target.id === "modal") modal.close();
});

clearLog();
refresh();