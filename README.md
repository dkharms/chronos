🕐 chronos is a continuous benchmarking tool that tracks performance metrics over time

![home-page](https://github.com/user-attachments/assets/558414f4-5a4b-44bd-b0fd-6f38126aba99)

Interactive version you can find [here](https://dkharms.github.io/chronos).

### Rationale

### Setting `chronos` up

#### Preconfiguration

Since `chronos` operates on two Git branches (by default these are `chronos-storage` and `gh-pages`)
you need create and push these branches to the `remote`:

```bash
git checkout --orphan chronos-storage \
  && git commit --allow-empty -m "[chronos] init" \
  && git push origin chronos-storage
```

And the same commands you must run for you GitHub Pages branch (if has not been created yet).

```bash
git checkout --orphan gh-pages \
  && git commit --allow-empty -m "[chronos] init" \
  && git push origin gh-pages
```

#### Action Configuration

Please refer to [this](https://github.com/dkharms/chronos/blob/main/.github/workflows/chronos.yml)
workflow file to grasp how to incorporate `chronos` into your GitHub Workflows.

And [here](https://github.com/dkharms/chronos/blob/main/action.yml) you find what inputs `chronos` is expecting.

Do not worry - it's dead simple! Just do not forget to pin the release version.

### Notes

I disclose that I've used AI assistance to build UI for data visualization (I am far from frontend-development),
but I've completed common-sense checks of the code that was spilt out by LLM - it's doable.
