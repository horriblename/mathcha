if vim.b.did_mathcha == 1 then
	return
end
vim.b.did_mathcha = 1

local err = require('mathcha').attach(vim.fn.bufnr())
if err then
	vim.notify_once(
		err,
		vim.log.levels.ERROR
	)
end
