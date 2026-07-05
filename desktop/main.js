const { app, BrowserWindow, dialog } = require('electron');
const childProcess = require('child_process');
const fs = require('fs');
const http = require('http');
const path = require('path');

const isPackaged = app.isPackaged;
const backendDir = isPackaged
  ? path.join(process.resourcesPath, 'backend')
  : path.join(__dirname, 'bin');
const backendExe = path.join(backendDir, process.platform === 'win32' ? 'aion.exe' : 'aion');

let backendProcess = null;
let mainWindow = null;
let appUrl = null;
let backendAddressFile = null;

function waitForBackendAddress(timeoutMs = 15000) {
  const startedAt = Date.now();

  return new Promise((resolve, reject) => {
    const probe = () => {
      fs.readFile(backendAddressFile, 'utf8', (error, address) => {
        if (!error && address.trim()) {
          resolve(`http://${address.trim()}`);
          return;
        }

        if (Date.now() - startedAt >= timeoutMs) {
          reject(new Error(`Backend did not report its address within ${timeoutMs}ms`));
          return;
        }
        setTimeout(probe, 100);
      });
    };

    probe();
  });
}

function waitForServer(url, timeoutMs = 15000) {
  const startedAt = Date.now();

  return new Promise((resolve, reject) => {
    const probe = () => {
      const request = http.get(url, (response) => {
        response.resume();
        resolve();
      });

      request.on('error', () => {
        if (Date.now() - startedAt >= timeoutMs) {
          reject(new Error(`Backend did not start within ${timeoutMs}ms`));
          return;
        }
        setTimeout(probe, 300);
      });

      request.setTimeout(1000, () => request.destroy());
    };

    probe();
  });
}

function startBackend() {
  if (!fs.existsSync(backendExe)) {
    throw new Error(`Backend executable not found: ${backendExe}`);
  }

  backendAddressFile = path.join(app.getPath('userData'), `backend-${process.pid}.addr`);
  fs.rmSync(backendAddressFile, { force: true });

  backendProcess = childProcess.spawn(backendExe, {
    cwd: backendDir,
    env: {
      ...process.env,
      AION_PORT_FILE: backendAddressFile,
    },
    windowsHide: true,
    stdio: 'ignore',
  });

  backendProcess.on('exit', () => {
    backendProcess = null;
  });
}

async function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1280,
    height: 820,
    minWidth: 980,
    minHeight: 640,
    show: false,
    autoHideMenuBar: true,
    webPreferences: {
      contextIsolation: true,
      nodeIntegration: false,
    },
  });

  mainWindow.once('ready-to-show', () => {
    mainWindow.show();
  });

  await mainWindow.loadURL(appUrl);
}

app.whenReady().then(async () => {
  try {
    startBackend();
    appUrl = await waitForBackendAddress();
    await waitForServer(appUrl);
    await createWindow();
  } catch (error) {
    dialog.showErrorBox('Battle Log 启动失败', error.message);
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

app.on('window-all-closed', () => {
  app.quit();
});

app.on('before-quit', () => {
  if (backendProcess) {
    backendProcess.kill();
  }
  if (backendAddressFile) {
    fs.rmSync(backendAddressFile, { force: true });
  }
});
