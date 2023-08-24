I want the tree structure output to look like the "tree" command would produce. E.g.

```
.
├── main.go
└── mytool
└── subdir/
          ... etc
```

And I want to the file name to print out before the file:

e.g.

>>>> subdir/thing.go <<<<

package subdir

// more code here
...

If there is a file called "project-requirements.md" that is a sibling to the directory where I am executing the script, I want to print that first.
