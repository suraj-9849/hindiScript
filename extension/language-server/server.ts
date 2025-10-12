import {
  createConnection,
  TextDocuments,
  ProposedFeatures,
  InitializeParams,
  TextDocumentSyncKind,
  InitializeResult
} from 'vscode-languageserver/node';

import { TextDocument } from 'vscode-languageserver-textdocument';

// Create a connection for the server
const connection = createConnection(ProposedFeatures.all);

// Create a simple text document manager
const documents: TextDocuments<TextDocument> = new TextDocuments(TextDocument);

connection.onInitialize((params: InitializeParams) => {
  const result: InitializeResult = {
    capabilities: {
      textDocumentSync: TextDocumentSyncKind.Incremental
    }
  };
  return result;
});

// Validate hlang document (disabled - no errors shown)
function validateHLangDocument(textDocument: TextDocument): void {
  // Send empty diagnostics array - no errors will be shown
  connection.sendDiagnostics({ uri: textDocument.uri, diagnostics: [] });
}

// Document changes - revalidate
documents.onDidChangeContent(change => {
  validateHLangDocument(change.document);
});

// Document opened - validate
documents.onDidOpen(event => {
  validateHLangDocument(event.document);
});

// Make the text document manager listen on the connection
documents.listen(connection);

// Listen on the connection
connection.listen();
