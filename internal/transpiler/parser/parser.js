#!/usr/bin/env node

const ts = require('typescript');
const fs = require('fs');
const path = require('path');

/**
 * Parse TypeScript file and output AST as JSON
 * Usage: node parser.js <input.ts>
 */
function parseTypeScriptFile(filePath) {
  const sourceCode = fs.readFileSync(filePath, 'utf8');
  
  const sourceFile = ts.createSourceFile(
    path.basename(filePath),
    sourceCode,
    ts.ScriptTarget.Latest,
    true
  );

  // Convert AST to a simplified JSON structure
  const simplifiedAST = convertNode(sourceFile);
  
  return JSON.stringify(simplifiedAST, null, 2);
}

/**
 * Convert TypeScript AST node to simplified JSON structure
 */
function convertNode(node) {
  const result = {
    kind: ts.SyntaxKind[node.kind],
    kindNumber: node.kind,
    pos: node.pos,
    end: node.end,
  };

  // Add text for identifiers and literals
  if (ts.isIdentifier(node)) {
    result.text = node.text;
  }
  
  if (ts.isStringLiteral(node) || ts.isNumericLiteral(node)) {
    result.text = node.text;
  }

  // Add text for template parts
  if (node.kind === ts.SyntaxKind.TemplateHead || 
      node.kind === ts.SyntaxKind.TemplateMiddle || 
      node.kind === ts.SyntaxKind.TemplateTail) {
    result.text = node.text;
  }

  // Add name for declarations
  if (node.name) {
    if (ts.isIdentifier(node.name)) {
      result.name = node.name.text;
    } else {
      // For binding patterns (destructuring), include the full pattern
      result.nameNode = convertNode(node.name);
    }
  }

  // Add operator for binary expressions
  if (node.operatorToken) {
    result.operator = ts.SyntaxKind[node.operatorToken.kind];
    result.operatorNumber = node.operatorToken.kind;
  }

  // Add questionDot for optional chaining
  if (node.questionDotToken) {
    result.questionDot = true;
  }

  // Add types for union types
  if (node.types) {
    result.types = node.types.map(convertNode);
  }

  // Add parameters for functions
  if (node.parameters) {
    result.parameters = node.parameters.map(convertNode);
  }

  // Add members for interfaces and classes
  if (node.members) {
    result.members = Array.from(node.members).map(convertNode);
  }

  // Add heritage clauses for class inheritance
  if (node.heritageClauses) {
    result.heritageClauses = Array.from(node.heritageClauses).map(convertNode);
  }

  // Add properties for object literals
  if (node.properties) {
    result.properties = Array.from(node.properties).map(convertNode);
  }

  // Add elements for arrays
  if (node.elements) {
    result.elements = Array.from(node.elements).map(convertNode);
  }

  // Add statements for blocks
  if (node.statements) {
    result.statements = Array.from(node.statements).map(convertNode);
  }

  // Add body for functions
  if (node.body) {
    result.body = convertNode(node.body);
  }

  // Add initializer for variables
  if (node.initializer) {
    result.initializer = convertNode(node.initializer);
  }

  // Add declarations for variable statements
  if (node.declarationList) {
    result.declarations = node.declarationList.declarations.map(convertNode);
  }

  // Add type for type annotations
  if (node.type) {
    result.type = convertNode(node.type);
  }

  // Add template expression properties
  if (node.head) {
    result.head = convertNode(node.head);
  }
  if (node.expression) {
    result.expression = convertNode(node.expression);
  }
  if (node.literal) {
    result.literal = convertNode(node.literal);
  }

  // Recursively convert children
  const children = [];
  ts.forEachChild(node, (child) => {
    // Skip already processed properties
    if (child !== node.name && child !== node.type && child !== node.body && 
        child !== node.initializer && child !== node.declarationList &&
        child !== node.head && child !== node.expression && child !== node.literal) {
      children.push(convertNode(child));
    }
  });

  if (children.length > 0) {
    result.children = children;
  }

  return result;
}

// Main execution
if (require.main === module) {
  const args = process.argv.slice(2);
  
  if (args.length === 0) {
    console.error('Usage: node parser.js <input.ts>');
    process.exit(1);
  }

  const inputFile = args[0];
  
  if (!fs.existsSync(inputFile)) {
    console.error(`File not found: ${inputFile}`);
    process.exit(1);
  }

  try {
    const ast = parseTypeScriptFile(inputFile);
    console.log(ast);
  } catch (error) {
    console.error('Error parsing TypeScript file:', error.message);
    process.exit(1);
  }
}

module.exports = { parseTypeScriptFile };
