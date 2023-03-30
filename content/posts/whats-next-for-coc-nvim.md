+++ 
date = 2023-03-30T21:17:45+01:00
title = "What's next for Coc.nvim"
description = "What's next for Coc.nvim"
slug = "whats-next-for-coc-nvim"
authors = []
tags = []
categories = []
externalLink = ""
series = []
+++


As a huge fan of Coc.nvim, I'm eager to see it being enjoyed by more people and
continuously improved, despite facing competition from built-in LSP clients
inside Neovim.

People may want "small" components to have composability and control over
things, which is what the built-in LSP client is about. For some users out
there, it will work, but for the vast majority of us, being built-in doesn't
automatically make it good. The more moving parts there are, the easier it is to
break. I believe people are not necessarily against the "all-in-one" solution,
as long as it _Just Works_ and still preserves the flexibility of Vim.

However, the more people use it, the faster it will improve. Coc.nvim should
find its niche and bet on it. One of such niche is the ability of porting VSCode
plugins to Coc.nvim, the accessibility to the vast (still growing) amount of
VSCode plugins is a huge upside. Enabling the "Easy Adaptation" of these plugins
would be a great feature for many users. While I see efforts being made in this
area, such as the pull request found here
(https://github.com/neoclide/coc.nvim/pull/3713), it has yet to be completed. In
an ideal world, wrapping an existing VSCode plugin to make it a Coc.nvim plugin
would be possible. While this may not be feasible in reality, it could be a
valuable direction to explore.

For those who are currently using VSCode, Coc.nvim should provide a seamless
transition into the Vim world without having to sacrifice the functionality and
features they've come to rely on. By offering accessibility to a vast array of
VSCode plugins and enabling easy adaptation of these plugins, Coc.nvim makes it
easier for VSCode users to switch to Vim while still being able to use the
plugins they're already familiar with. This way, they can enjoy the speed and
efficiency of Vim without giving up the tools they rely on.

For those who are already familiar with Vim, Coc.nvim is well-positioned to
offer them a reliable and efficient LSP client. In addition, they can tap into
the VSCode world, especially considering that many of the AI plugins are
released on VSCode first. You don't want to miss the AI train, right?

In summary, Coc.nvim is attracting both VSCode users looking to improve their
editing skills by learning Vim, as well as Vim users seeking to tap into the
VSCode world.

Also posted on https://github.com/neoclide/coc.nvim/discussions/4593
