package network

import (
    "net/http"
)

func MakeMethodChecker( mflag int ) map[string] bool {

    fl := make( map[string] bool )
    fl[http.MethodGet]     = MethodGet     & mflag == MethodGet
    fl[http.MethodHead]    = MethodHead    & mflag == MethodHead
    fl[http.MethodPost]    = MethodPost    & mflag == MethodPost
    fl[http.MethodPut]     = MethodPut     & mflag == MethodPut
    fl[http.MethodPatch]   = MethodPatch   & mflag == MethodPatch
    fl[http.MethodDelete]  = MethodDelete  & mflag == MethodDelete
    fl[http.MethodConnect] = MethodConnect & mflag == MethodConnect
    fl[http.MethodOptions] = MethodOptions & mflag == MethodOptions
    fl[http.MethodTrace]   = MethodTrace   & mflag == MethodTrace

    return fl

}