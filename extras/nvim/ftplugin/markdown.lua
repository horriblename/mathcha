if vim.g.skip_mathcha or vim.b.did_mathcha == 1 then
	return
end
vim.b.did_mathcha = 1

-- extremely buggy if attaching before treesitter starts so
-- here we are :))))
vim.schedule(function()
	local err = require('mathcha').attach(vim.fn.bufnr())
	if err then
		vim.notify_once(
			err,
			vim.log.levels.ERROR
		)
	end
end)
