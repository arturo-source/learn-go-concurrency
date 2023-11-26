# Exercise: Web Crawler

> This exercise is taken from <https://go.dev/blog/pipelines#digesting-a-tree> and modified to be a little bit more difficult. You can find my solution in [this respository](https://github.com/arturo-source/repeated-files-detector).

Use Go's concurrency to make a program that receives a folder and evaluate which files are duplicated. The user can indicate if wants the output into a file, otherwise the output will be written in stdout. You have to use MD5 algorithm for hashing, instead of comparing bytes.
