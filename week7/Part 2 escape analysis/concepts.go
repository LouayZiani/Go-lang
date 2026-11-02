// escape analysis is the part of the compiler that decides whether the value stays in the fct's stack frame,
// or might it be used AFTER the function returns. if it might, it ESCAPES to the heap

// Inlining copies small fcts into their callers --> can change escape decisions

// does not escape: stays on the stack
/* what makes a value escape?
Sharing up: returning a ptr to a local (return &u)
Storing a pointer somehwere that lives longer: a globa, map, slice that escapes, a channel
Interfaces: passing a value as any/interface usually boxes it on the heap (boxing)
Closures that outlive the function capture their variables on the heap
Size: too big for the stack
*/

package main
