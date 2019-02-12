package sprout

/*
 + SPACE FACTORY
 *
 * A pseudo-static member that manages space creation
 */

type space_factory struct {}
var  static_space_factory = &space_factory{}

func SpaceFactory() *space_factory {}
func( _spcfac *space_factory ) New( _name string ) *Space {}