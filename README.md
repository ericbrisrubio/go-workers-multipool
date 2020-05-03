# go-workers-multipool
Manager for multiple pools of workers in Golang.

This is a wrapper to manage multiple pool of workers of the kind:
https://github.com/enriquebris/goworkerpool

You can easily have many pool to execute different tasks and manage them separately.

### Example Use Case:
System to resize images

The system could use 2 different pools of workers, one to process low size images (< 5 MB), and a second pool to process 
the rest (> 5 MB). Since per the history, 70% of elements coming are bigger than 5 MB, and we want to do some extra work 
on them without affecting the conversion of the smaller ones, the pool for the bigger images could have more workers than 
the other. If at some point, we determine that we need only half of the current amount of workers converting big images, 
we can cut them in half on the fly without affecting anything else and without restarting the system.

Current Features:
- Add multiple pools to be managed at the same time
- Define a function for the workers on a specific pool
- Start the workers on a pool
- Edit the amount of workers for a pool on the fly
- Kill multiple workers on a pool on the fly
- Add workers for a pool on the fly
- Pause all the workers for a pool
- Resume all the workers for a pool

System Overview:

![system-overview](./go-workers-multipool-Overview.svg)

### Quick Start: 
(Based on the previous use case about images conversion)


