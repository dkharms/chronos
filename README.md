### Setting `chronos` up

Since `chronos` operates on two Git branches (by default these are `chronos-storage` and `gh-pages`)
you need create and push these branches to the `remote`:

```bash
git checkout --orphan chronos-storage \
  && git commit --allow-empty -m "[chronos] init" \
  && git push origin chronos-storage
```

And the same commands you must run for you GitHub Pages branch (if has not been created yet).
