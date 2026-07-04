package gsxgrammar

import _ "embed"

//go:embed queries/gsx/highlights.scm
var GSXHighlights []byte

//go:embed queries/gsx/injections.scm
var GSXInjections []byte

//go:embed queries/go/highlights.scm
var GoHighlights []byte

//go:embed queries/javascript/highlights.scm
var JSHighlights []byte

//go:embed queries/css/highlights.scm
var CSSHighlights []byte
