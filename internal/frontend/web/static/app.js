'use strict';

// --- Tab switching ---
document.querySelectorAll('.tab-btn').forEach(btn => {
  btn.addEventListener('click', () => {
    document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    btn.classList.add('active');
    document.getElementById('tab-' + btn.dataset.tab).classList.add('active');
  });
});

// --- API helpers ---
async function apiFetch(path, options = {}) {
  const res = await fetch(path, {
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.message || `Ошибка ${res.status}`);
  return body;
}

function showMsg(el, text, isOk) {
  el.textContent = text;
  el.className = 'form-msg ' + (isOk ? 'ok' : 'err');
}

// --- Currencies tab ---
async function loadCurrencies() {
  const el = document.getElementById('currencies-list');
  try {
    const data = await apiFetch('/currencies');
    if (!data.length) { el.innerHTML = '<p class="loading">Валюты не найдены</p>'; return; }
    el.innerHTML = `
      <div class="table-wrap">
        <table>
          <thead><tr><th>Код</th><th>Название</th><th>Символ</th></tr></thead>
          <tbody>${data.map(c => `
            <tr>
              <td><strong>${esc(c.code)}</strong></td>
              <td>${esc(c.name)}</td>
              <td>${esc(c.sign)}</td>
            </tr>`).join('')}
          </tbody>
        </table>
      </div>`;
  } catch (e) {
    el.innerHTML = `<p class="error-msg">${esc(e.message)}</p>`;
  }
}

// --- Rates tab ---
async function loadRates() {
  const el = document.getElementById('rates-list');
  try {
    const data = await apiFetch('/exchangeRates');
    if (!data.length) { el.innerHTML = '<p class="loading">Курсы не найдены</p>'; return; }
    el.innerHTML = `
      <div class="table-wrap">
        <table>
          <thead><tr><th>Пара</th><th>Базовая</th><th>Целевая</th><th>Курс</th></tr></thead>
          <tbody>${data.map(r => `
            <tr>
              <td><strong>${esc(r.baseCurrency.code)}/${esc(r.targetCurrency.code)}</strong></td>
              <td>${esc(r.baseCurrency.name)}</td>
              <td>${esc(r.targetCurrency.name)}</td>
              <td>${r.rate.toFixed(6)}</td>
            </tr>`).join('')}
          </tbody>
        </table>
      </div>`;
  } catch (e) {
    el.innerHTML = `<p class="error-msg">${esc(e.message)}</p>`;
  }
}

// --- Converter ---
document.getElementById('converter-form').addEventListener('submit', async e => {
  e.preventDefault();
  const from = document.getElementById('conv-from').value.trim().toUpperCase();
  const to = document.getElementById('conv-to').value.trim().toUpperCase();
  const amount = document.getElementById('conv-amount').value;
  const el = document.getElementById('converter-result');

  if (from.length !== 3 || to.length !== 3) {
    el.innerHTML = '<p class="error-msg">Введите корректные коды валют (3 буквы)</p>';
    return;
  }

  el.innerHTML = '<p class="loading">Вычисление...</p>';
  try {
    const r = await apiFetch(`/exchange?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}&amount=${encodeURIComponent(amount)}`);
    el.innerHTML = `
      <div class="result-card">
        <div class="row"><span>Из</span><span>${esc(r.baseCurrency.name)} (${esc(r.baseCurrency.sign)})</span></div>
        <div class="row"><span>В</span><span>${esc(r.targetCurrency.name)} (${esc(r.targetCurrency.sign)})</span></div>
        <div class="row"><span>Курс</span><span>1 ${esc(r.baseCurrency.code)} = ${r.rate.toFixed(6)} ${esc(r.targetCurrency.code)}</span></div>
        <div class="row"><span>Сумма</span><span>${r.amount} ${esc(r.baseCurrency.code)}</span></div>
        <div class="result-amount">${r.convertedAmount.toFixed(2)} ${esc(r.targetCurrency.code)}</div>
      </div>`;
  } catch (e) {
    el.innerHTML = `<p class="error-msg">${esc(e.message)}</p>`;
  }
});

// --- Add currency ---
document.getElementById('add-currency-form').addEventListener('submit', async e => {
  e.preventDefault();
  const msg = document.getElementById('currency-msg');
  const fd = new FormData(e.target);
  const body = {
    code: fd.get('code').toUpperCase(),
    name: fd.get('name'),
    sign: fd.get('sign'),
  };

  try {
    await apiFetch('/currencies', { method: 'POST', body: JSON.stringify(body) });
    showMsg(msg, `Валюта ${body.code} успешно добавлена`, true);
    e.target.reset();
    loadCurrencies();
  } catch (err) {
    showMsg(msg, err.message, false);
  }
});

// --- Add/Update rate ---
document.getElementById('rate-form').addEventListener('submit', async e => {
  e.preventDefault();
  const msg = document.getElementById('rate-msg');
  const fd = new FormData(e.target);
  const action = e.submitter?.value ?? 'create';
  const base = fd.get('baseCurrencyCode').toUpperCase();
  const target = fd.get('targetCurrencyCode').toUpperCase();
  const rate = parseFloat(fd.get('rate'));

  try {
    if (action === 'update') {
      await apiFetch(`/exchangeRate/${base}${target}`, {
        method: 'PATCH',
        body: JSON.stringify({ rate }),
      });
      showMsg(msg, `Курс ${base}/${target} обновлён`, true);
    } else {
      await apiFetch('/exchangeRates', {
        method: 'POST',
        body: JSON.stringify({ baseCurrencyCode: base, targetCurrencyCode: target, rate }),
      });
      showMsg(msg, `Курс ${base}/${target} добавлен`, true);
      e.target.reset();
    }
    loadRates();
  } catch (err) {
    showMsg(msg, err.message, false);
  }
});

// --- XSS guard ---
function esc(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

// --- Init ---
loadCurrencies();
loadRates();
