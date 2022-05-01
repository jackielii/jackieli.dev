+++
date = 2022-05-01T22:11:00+01:00
title = "Vim jump to caps word on the same line"
description = "Jump to caps word in the same line in VIM/Neovim"
slug = "vim-jump-to-caps"
tags = ["vim", "neovim"]
categories = ["vim"]
externalLink = ""
series = []
+++

I often find myself wanting to jump to the word starting with capital letter on
the same line. Especially in Go where the methods you import package always
starts with capital letters. E.g.

```
<> fmt.Println("Hello World")
```

When you cursor is at `<>` there are a few ways to quickly jump to `Println`:

1. press `w` **three** times (because there is a dot)
2. press `fP`
3. press `l` many times until you're at `P`

But it gets tricker when you have longer lines like this one:

```
<>	log := logging.NewLogrusLoggerLevel(logrus.InfoLevel)
```

If you want to jump to `I` of `InfoLevel`, you probably jump to end with `$`
then `b` to back one word. This is ok, but you have to think for a second: do I
press `w` and hold until I'm there? Or do I jump to the end and move back?

I find myself using this one so often and it got really annoying. So I made
this mapping:

```
" L/H to move to the next/previous capital letter of a word
nmap <silent> L :call search('\<\u', '', line('.'))<CR>
nmap <silent> H :call search('\<\u', 'b', line('.'))<CR>
```

It uses `search` so it doesn't add anything to your search register. I used `H`
& `L`. It feels similiar to move left and right `h` & `l`, but faster.

It's especially convenient because when I use this jump, I also want to query
the documentation of such methods using `K`:

```
nnoremap <silent> K :call <SID>show_documentation()<CR>

function! s:show_documentation()
  if (index(['vim','help'], &filetype) >= 0)
    execute 'h '.expand('<cword>')
  elseif (coc#rpc#ready())
    call CocActionAsync('doHover')
  else
    execute '!' . &keywordprg . " " . expand('<cword>')
  endif
endfunction
```

So often I only have to press `H` followed by `K` - _only three_ key presses
including `Shift` when I moved to a line and wanting to quickly jump to a
public exported method and query it's documentation.
