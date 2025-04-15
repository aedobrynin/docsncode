import * as vscode from 'vscode';
import { exec } from 'child_process';
import * as path from 'path';

export function activate(context: vscode.ExtensionContext) {
	const disposable = vscode.commands.registerCommand('docsncode-vscode-extension.runDocsncode', () => {
		const editor = vscode.window.activeTextEditor;

		if (!editor) {
			vscode.window.showInformationMessage('There is no current file');
			return;
		}


		const document = editor.document;
		const filePath = document.fileName;

		const config = vscode.workspace.getConfiguration('docsncode');
		const docsncodePath = config.get<string>('path', 'docncode');

		exec(`${docsncodePath} ${filePath}`, (error, stdout, stderr) => {
			if (error) {
				vscode.window.showErrorMessage(`docsncode error: ${stderr}`);
				return;
			}

			const panel = vscode.window.createWebviewPanel(
				'docsncode',
				'Docsncode',
				vscode.ViewColumn.Two,
				{}
			);

			panel.webview.html = stdout;
		});
	});

	context.subscriptions.push(disposable);
}

export function deactivate() { }
