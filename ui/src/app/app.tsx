
import styles from "./app.module.css";

import { AuthType, createClient, FileStat, ResponseDataDetailed } from "webdav";
import React from "react";
import ProgressBar from "components/progress_bar";

const client = createClient("http://localhost:7086/webdav", {
  authType: AuthType.Password,
  username: "flydav",
  password: "pass-for-testing" // sha256: 8c024a1ccc39abc26d05bacc4ab64b78ad4c4378e67df31975d5303ca2258a22
});

const dirname = (path: String) => {
  let ret = path.replace(/\\/g, '/').replace(/\/[^/]*$/, '');
  if (ret.length == 0) {
    return "/"
  }
  return ret;
}

const dateFormat = (input: string) => {
  const date = new Date(input)
  return date.toJSON();

}

interface FileStatExtended extends FileStat {
  isParent: boolean
}

const App = (): JSX.Element => {

  const [path, setPath] = React.useState<string>("/");
  const [pathUncommitted, setPathUncommitted] = React.useState<string>("/");
  const [files, setFiles] = React.useState<FileStatExtended[]>([]);
  const [loading, setLoading] = React.useState<boolean>(false);
  const [downloadProgress, setDownloadProgress] = React.useState(0);
  React.useEffect(() => {
    setLoading(true);
    pathUncommitted != path && setPathUncommitted(path);
    client.getDirectoryContents(path).then((files: FileStat[] | ResponseDataDetailed<FileStat[]>) => {
      // if is ResponseDataDetailed, unwrap
      if (!Array.isArray(files)) {
        alert("Error: unexpected response type")
      }
      let filesUnwrapped = files as FileStatExtended[];

      filesUnwrapped.map(f => {
        f.isParent = false
      })

      setFiles(filesUnwrapped);
      setLoading(false);
    });
  }, [path]);

  async function handleClickFile(file: FileStat) {
    const buff: Buffer = await client.getFileContents(file.filename, {
      format: "binary",
      onDownloadProgress: e => {
        setDownloadProgress(e.loaded / e.total * 100)
      },
    }) as Buffer
    saveBuffer(buff, file.basename)

  }
  const saveBuffer = (buf: Buffer, filename: string) => {
    const a = document.createElement('a');
    a.style.display = 'none';
    document.body.appendChild(a);
    const blob = new Blob([buf], { type: 'octet/stream' });
    const url = window.URL.createObjectURL(blob);
    a.href = url;
    a.download = filename;
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  };
  return (
    <main className={styles.main}>
      <header className="flex justify-between items-center p-4 border-b border-gray-300">
        <h3 className="text-2xl font-light">
          FlyDav UI
        </h3>
        <div className="flex">
          <input type="text" className={styles.searchFileInput} placeholder="Search files">
          </input>
          <button className={styles.buttonDefault} >Settings</button>
        </div>
      </header>

      {downloadProgress > 0 && downloadProgress < 100 && <ProgressBar current={downloadProgress} total={100}></ProgressBar>}
      <section className="p-4">
        <div className="flex items-center">
          <button className={styles.buttonDefault} onClick={() => setPath(dirname(path))}>To parent</button>
          <input type="text" className={styles.pathInput} value={pathUncommitted} onChange={(e) => setPathUncommitted(e.target.value)}></input>
          <button className={styles.buttonPrimary} onClick={() => setPath("/")}>Go</button>
        </div>
      </section>
      <section className="p-4">
        <div className={styles.fileList}>
          <div className={styles.fileListHeader}>
            <div className={styles.fileListCell}>Name</div>
            <div className={styles.fileListSizeCell}>Size</div>
            <div className={styles.fileListLastmodCell}>Last Modified</div>
          </div>

          {
            files.filter(file => file.type == "directory").map((file: FileStatExtended, idx, arr) => {
              return (
                <div className={styles.fileListItem}>
                  <div className={styles.fileListCell}>
                    <svg className={styles.fileListIcon} xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="#e79f2d" viewBox="0 0 16 16">
                      <path d="M9.828 3h3.982a2 2 0 0 1 1.992 2.181l-.637 7A2 2 0 0 1 13.174 14H2.825a2 2 0 0 1-1.991-1.819l-.637-7a1.99 1.99 0 0 1 .342-1.31L.5 3a2 2 0 0 1 2-2h3.672a2 2 0 0 1 1.414.586l.828.828A2 2 0 0 0 9.828 3zm-8.322.12C1.72 3.042 1.95 3 2.19 3h5.396l-.707-.707A1 1 0 0 0 6.172 2H2.5a1 1 0 0 0-1 .981l.006.139z" />
                    </svg>
                    <a className={styles.entryLink}
                      onClick={() => setPath(file.filename)}
                    >{file.basename}/</a></div>
                  <div className={styles.fileListSizeCell}>{file.size}</div>
                  <div className={styles.fileListLastmodCell}>{dateFormat(file.lastmod)}</div>
                </div>
              )
            })
          }
          {
            files.filter(file => file.type == "file").map((file: FileStat, idx, arr) => {
              return (
                <div className={styles.fileListItem}>
                  <div className={styles.fileListCell}>
                    <svg className={styles.fileListIcon} xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M4 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H4zm0 1h8a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1z" />
                    </svg>
                    <a onClick={async () => { await handleClickFile(file) }} className={styles.entryLink}>{file.basename}</a></div>
                  <div className={styles.fileListSizeCell}>{file.size}</div>
                  <div className={styles.fileListLastmodCell}>{dateFormat(file.lastmod)}</div>
                </div>
              )
            })
          }
        </div>
      </section>
      <footer className={styles.footer}>
        <div className={styles.copyRight}>
          Copyright {new Date().getFullYear()} - <a className={styles.projectLink} href="https://github.com/pluveto/flydav">FlyDav</a> Project
        </div>
      </footer>
    </main>
  );
};

export default App;
