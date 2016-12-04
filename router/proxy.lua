local router = require 'router'
local json = require "cjson.safe"
local r = router.new()

function split(inputstr, sep)
  local t, i = {}, 1
  if null == inputstr then
    return t
  end
  for str in string.gmatch(inputstr, "([^"..sep.."]+)") do
    t[i] = str
    i = i + 1
  end
  return t
end

function ok(data)
  return {
    status = 200,
    data = data,
    error = ''
  }
end

function error(message, code)
  return {
    status = code,
    data = {},
    error = message
  }
end

r:get('/api/v1/jobs_queue', function(params)
  res = ngx.location.capture("/airtableJob")
  res_data = (json.decode(res.body))
  if res.status ~= 200 then
    return error(res_data.error.message, 400)
  else
    resp = {}
    for i, record in ipairs(res_data.records) do
      resp[i] = record.fields
    end
    return ok(resp)
  end
end)

r:post('/api/v1/jobs_queue', function(params)

  local req, err = json.decode(ngx.req.get_body_data())

  if req == nil then req = {} end

  req.status = "new"

  for k, field in ipairs ({"author", "kind", "callback", "params"}) do
    if nil == req[field] then
      return error("'" .. field .. "' field is required", 400)
    end
  end

  res = ngx.location.capture("/airtableJob", {
    method = ngx.HTTP_POST,
    body = json.encode({fields = req}),
  })

  res_data = (json.decode(res.body))
  if res.status ~= 200 then
    return error(res_data.error.message, 400)
  else
    return ok({id = res_data.id})
  end

end)

r:put('/api/v1/jobs_queue/:id', function(params)

  local req, err = json.decode(ngx.req.get_body_data())

  if req == nil then req = {} end

  for k, field in ipairs ({"author", "kind", "callback", "params"}) do
    if nil == req[field] then
      return error("'" .. field .. "' field is required", 400)
    end
  end

  res = ngx.location.capture("/airtableJob/" .. params.id, {
      method = ngx.HTTP_PUT,
      body = json.encode({ fields = req })
  })

  res_data = (json.decode(res.body))
  if res.status ~= 200 then
    return error(res_data.error.message, 400)
  else
    return ok({id = params.id})
  end

end)

r:get('/api/v1/status', function(params)
  return ok({
    name =  "Korchagin Stanislav",
    role = "Team Lead / Tech Lead / CTO",
    category = "IT",
    icon = "https://lh3.googleusercontent.com/O6pzfqAmTVKoAjmiqzK-hd5hZ4xTKnyHZoeCUplAT1D2SQZVMovlLOe8uTxJKtlSLO7Bsx0Mw8cc=w665-h1000-no",
    description = [[
      Успешно разрабатывал две(или три) социальные сети, пару CRM-систем, одну BPM, пару платежных витрин с генератором, генератор магазинов, несколько поисковых систем для больших порталов, организовывал отдел разработок, строил релиз-циклы и процессы разработки. А, еще один маленький CDN. Переставать программировать не готов.

      - Управление программистами, требованиями и процессами
      - Много php, немного java, python, lua, go, etc
      - Design patterns, OOP, TDD, KISS, YAGNI и другие страшные аббревиатуры
      - Unit-тесты, приемочное тестирование, стресс тестирование, проверка серверов на пылестойкость и ударопрочность
      - Извращенные архитектуры, многопоточность, межмашинковое взаимодействие, api, события, масштабирование вширь и вглубь, алгоритмы под данные, данные под задачи, очереди, map/reduce, hash tables, нарезка базы кубиками и колечками
      - Извращенные архитектуры, многопоточность, межмашинковое взаимодействие, алгоритмы под данные, данные под задачи, очереди, map/reduce, hash tables, нарезка базы кубиками и колечками
      - Данные больше 100Gb и хиты больше 10M
      ]],
    location = {
      city = "Kiev",
      countryCode = "UA"
    },
    links = {
      {
        network = "Twitter",
        username = "korchasa",
        url = "https://twitter.com/korchasa"
      },
      {
        network = "Github",
        username = "korchasa",
        url = "https://github.com/korchasa"
      }
    },
    locales = {
      {
        language = "Russian",
        fluency = "Native speaker"
      },
      {
        language = "English",
        fluency = "Чтение документации, переписка, называние переменных"
      }
    }
  })
end)

ngx.header.content_type = 'application/json';
local isRouted, result = r:execute(
  ngx.var.request_method,
  ngx.var.request_uri,
  ngx.req.get_uri_args(),
  ngx.req.get_post_args()
)

if true ~= isRouted then
  result = error(result, 404)
end

ngx.status = result.status

if 400 > result.status then
  ngx.print(json.encode(result))
else
  ngx.print(json.encode({
    status = result.status,
    message = result.error
  }))
  ngx.log(ngx.ERR, result.error)
end
