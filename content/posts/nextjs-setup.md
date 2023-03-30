+++
draft = true
date = 2022-09-23T17:41:45+04:00
title = ""
description = ""
slug = ""
authors = []
tags = []
categories = []
externalLink = ""
series = []
+++

```bash
yarn create next-app --typescript my-app
yarn add --dev @typescript-eslint/eslint-plugin
yarn add --dev prettier eslint-config-prettier
```

`.eslintrc.json`

```json
{
  "extends": [
    "next/core-web-vitals",
    "plugin:@typescript-eslint/recommended",
    "prettier"
  ]
}
```

`.prettierrc.json`

```json
{
  "semi": false,
  "trailingComma": "es5",
  "singleQuote": true,
  "tabWidth": 2,
  "useTabs": false
}
```

```bash
yarn add --dev husky lint-staged
yarn husky install
yarn husky add .husky/pre-commit "npx lint-staged"
```

`lint-staged.config.js`

```javascript
module.exports = {
  // Type check TypeScript files
  '**/*.(ts|tsx)': () => 'yarn tsc --noEmit',

  // Lint then format TypeScript and JavaScript files
  '**/*.(ts|tsx|js)': (filenames) => [
    `yarn eslint --fix ${filenames.join(' ')}`,
    `yarn prettier --write ${filenames.join(' ')}`,
  ],

  // Format MarkDown and JSON
  '**/*.(md|json)': (filenames) =>
    `yarn prettier --write ${filenames.join(' ')}`,
}
```
