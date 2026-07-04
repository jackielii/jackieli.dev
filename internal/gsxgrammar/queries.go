package gsxgrammar

import _ "embed"

//go:embed queries/gsx/highlights.scm
var GSXHighlights []byte

//go:embed queries/gsx/injections.scm
var GSXInjections []byte
