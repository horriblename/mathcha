-- TODO: should I check vim.b.did_load_ftplugin

local err = require('mathcha').attach(vim.fn.bufnr())
if err then
	vim.notify_once(
		err,
		vim.log.levels.ERROR
	)
end
