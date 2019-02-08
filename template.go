package sprout

import (
    "crpyto/sha256"
    "fmt"
    "html/template"
)

/*
 # Realtime Template Specification
 
 - A realtime template must have a unique name, which is hash of the body
 */

const (
    template_token_nil = -1 + iota
    template_token_left_delimiter
    template_token_right_delimiter
    template_token_string
    template_token_space // template's lexer treats ' ' and '\t' as spaces
)

type template_tree struct {
    body string
    curr *template_token
}

type template_token struct {
    typ int
    val string
    pos int
}

func make_template_tree( _body string ) ( *template_tree ) {
    return &template_tree{
        body: _body,
        curr: nil,
    }
}

func ( tr *template_tree ) next() ( *template_token ) {
    
    _piece := ""
    _last_type := template_token_nil
    _pos := tr.curr.pos + len( tr.curr.val )
    
    if _pos >= len( tr.body ) {
        return nil
    }
    
    _body := tr.body[_pos:]
    
    for i, _rune := range _body {
        
        _curr_piece := _rune
        _curr_type  := template_token_nil
        
        switch _rune {
        case ' ': fallthrough
        case '\t':
            _curr_type = template_token_space
        case template_left_delimiter[0]:
            if i + len( template_left_delimiter ) - 1 >= len( _body ) {
                _curr_type = template_token_string
            }
        }
        
        
        
    }
    
}

func ( sp *Sprout ) recursive_parse_realtime_template( _body string ) ( *template.Template, error ) {
    
    // Prepare Name
    _sha256 := sha256.New()
    _sha256.Write( []byte( _body ) )
    _name := fmt.Sprintf( "%x", _sha256.Sum( nil ) )
    
    
    
}