// MCP (Model Context Protocol) service client
// This provides client-side functionality for MCP server interactions

export interface MCPManifest {
  name: string;
  description?: string;
  version?: string;
  tools?: any[];
  resources?: any[];
  prompts?: any[];
}

export interface MCPInstallationCheck {
  installed: boolean;
  version?: string;
  error?: string;
}

export class MCPServiceClient {
  async checkInstallation(data: any, signal?: AbortSignal): Promise<MCPInstallationCheck> {
    console.warn('MCPServiceClient.checkInstallation not implemented', data);
    // Return a basic response for build purposes
    return {
      installed: false,
      error: 'MCP service not fully implemented',
    };
  }

  async getStdioMcpServerManifest(data: any): Promise<MCPManifest> {
    console.warn('MCPServiceClient.getStdioMcpServerManifest not implemented', data);
    throw new Error('MCP service not fully implemented');
  }

  async getStreamableMcpServerManifest(data: any): Promise<MCPManifest> {
    console.warn('MCPServiceClient.getStreamableMcpServerManifest not implemented', data);
    throw new Error('MCP service not fully implemented');
  }

  async listTools(data: any): Promise<any[]> {
    console.warn('MCPServiceClient.listTools not implemented', data);
    return [];
  }

  async listResources(data: any): Promise<any[]> {
    console.warn('MCPServiceClient.listResources not implemented', data);
    return [];
  }

  async listPrompts(data: any): Promise<any[]> {
    console.warn('MCPServiceClient.listPrompts not implemented', data);
    return [];
  }

  async callTool(params: any, toolName: string, args: any): Promise<any> {
    console.warn('MCPServiceClient.callTool not implemented', params, toolName, args);
    throw new Error('MCP service not fully implemented');
  }
}

export const mcpService = new MCPServiceClient();
