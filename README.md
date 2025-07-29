# Et tu, CESR? <!-- the sharpest parser in the Senate ⚔️ -->

**Et tu, CESR?** is a tiny Go CLI that “stabs” a CESR stream—
the self‑framing event format used by KERI & ACDC—and spits out
pretty‑printed JSON bodies plus the size of each attachment block.

> *“An event a day keeps the Brutus away.”*

---

## ✨ Features

| What                | Why it’s handy                                                                                          |
| ------------------- | ------------------------------------------------------------------------------------------------------- |
| **Self‑contained**  | Pure Go 1.22+, no external libraries.                                                                   |
| **Multi‑event**     | Handles huge concatenated streams effortlessly.                                                         |
| **Readable output** | Pretty JSON with event \`t\`ype and \`sn\` sequence.                                                    |
| **Straightforward** | A single regex on \`{\\"v\\":\\"KERI10…\` / \`{\\"v\\":\\"ACDC10…\`, header gives body length—no magic. |

---

## 🔧 Installation

```bash
git clone https://github.com/ForkbombEu/et-tu-cesr.git
cd et-tu-cesr
mise i
task
```

Or run once without building:

```bash
task run
```

---

## 🚀 Quick start

```bash
./et-tu-cesr samples/E4OU1DuxIAtRRscHSSQCO0UIpk3tVc0QHaNBDUmpHKac-acdc.cesr
```

```text
### Event 1  (t=icp  sn=0)
{
  "v": "KERI10JSON0000fb_",
  "t": "icp",
  "d": "EL1L56Lyo…Bug",
  ...
}
• attachment bytes: 160

### Event 2  (t=ixn  sn=1)
{ … }
• attachment bytes: 188
```

---

## 🏗️ How it works (bird’s‑eye)

1. **Scan** for \`{\\"v\\":\\"KERI…\` / \`{\\"v\\":\\"ACDC…\` with a regex.
2. **Parse** the 17‑char CESR header (e.g. \`KERI10JSON000249\_\`) to get body length.
3. **Slice** \`bodyStart : bodyEnd\`, \`json.Unmarshal\` → pretty print.
4. **Attachments** = \`bodyEnd : nextHeader\` → just measure the bytes.
5. Repeat until no more headers.

---

## 🛣️ Roadmap

* [ ] `--raw` flag to dump attachment hex
* [ ] Stream from **stdin**
* [ ] ANSI color highlight for key fields (\`t\`, \`d\`, \`i\`)
* [ ] Unit tests against official KERI/ACDC samples

---

## 🤝 Contributing

Pull requests welcome!

1. Fork this repo
2. \`git checkout -b feature/your-idea\`
3. Commit + \`go vet ./...\`
4. Open a PR and describe what changed

---

## 📜 License

MIT License — see \`LICENSE\` for full text.

---

> *“Et tu, CESR?” — because even the most cryptic files can be betrayed and spill their secrets.*

