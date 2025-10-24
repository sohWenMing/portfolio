// The dbinteface package is meant to allow the choice of database that is used in the main
// application, or test instances to be interchangeable. To use the interface, all types
// defined within the dbinterface module should be **IMPORTED** into the actual wrapper around
// the database functionality.
//
// Within the database wrapper module, define inputs and outputs using the types defined and
// imported from the dbinterface module, wrapping the actual methods/functions exposed from
// the database layer, ensuring that the resulting wrapping function fulfils the interface
// defined in the interface module.
//
// End user comsumers (like handlers in the main module) should be defined to use the
// **INTERFACES** in the dbinterface module, and not the actual functions from the databsase
// layer
//

package dbinterface
