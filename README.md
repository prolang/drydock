# Drydock
Drydock is a cooperative multitasking framework for Go.

## Goals
In some scenarios it is both convenient and efficient to express cooperative concurrency instead of
true parallelism (e.g. goroutines).  Drydock was built to provide:

1.  **Very fine-grained concurrency:**  Asynchronous computations, async results, and turns are VERY
cheap.  It is appropriate to create a very large number of them to express very fine-grained
concurrency.

2.  **Reduced Interleavings:**  Cooperative concurrency gives code direct control over when it
yields to other concurrent code.  This allows code to know deterministically when the state it is
working with has **NOT** been changed by another concurrent party.  So called TOCTOU bugs.  This
significantly reduces the number of possible interleavings (i.e. places where one thread yields
control to another thread which executes for a while before returning control to the original
thread, possibly after mutating some state) that a program may encounter. Reduced interleavings make
it easier to reason about program behavior, to model that behavior (e.g. through state machines),
and to log that behavior at runtime for debugging purposes.

3.  **Single-threaded containers:**  All async computations execute within single-threaded
containers called Actors (in the Erlang sense).  Actors only communicate with each other through a
strict linear (ownership transfer) message passing paradigm.  Since Actors are ALWAYS singled-
threaded, no code running in a Actor ever needs to use locks or other synchronization primitives
(e.g. sync.Mutex, sync.Cond, atomic.*, etc.).  The incorrect use of locks is one of the most common
bugs in multi-threaded programming.  The difficulty in both debugging and fixing these kinds of
errors is what makes multi-threaded programming HARD.  Additionally, the absence of locking can
noticeably increase performance in some scenarios.

4.  **Deterministic Completion:**  All async computation return a completion result (e.g. async.R).
A completion result is a placeholder that becomes resolved when the asynchronous computation has
finished whether it finished successfully, encountered an error, or was cancelled.  All asynchronous
computation return a completion result even if they don't have a return value (i.e. all async
function return at least an error that may be nil on success).  This guarantees that asynchronous
computations can be composed to form other asynchronous computation which are similarly guaranteed
to complete.

## Work In Progress

Drydock is still a work in progress.  There still remain many elements of the framework that are
still under development or not yet implemented.

Here is a brief list of WIP tasks (in no particular order):

* Sample code.  The framework needs sample code that demonstrates how to use it.

* Cancellation.  Most asynchronous computation is subject to cancellation and a composable 
cancellation framework is needed that both encourages efficient termination of unneeded computation 
and guarantees computation is not orphaned.

* Code generation for strongly-typed result subtypes.  E.g. async.R (void),  async.StringR (string).
There should be a generated type for any value that meets might be the payload of a completion
result including all language primitives and any complex types that meet the linear requirements.

* Formalization of the interface definition for asynchronous message passing, whether between 
concurrency domain within the same actor or across actors.
