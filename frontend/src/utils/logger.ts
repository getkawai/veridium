/**
 * Logger utility for debugging database models
 * Uses Wails log service for consistent logging across the app
 */
import { LogService, Level } from '@@/github.com/wailsapp/wails/v3/pkg/services/log';

/**
 * Log levels for easy reference
 */
export const LogLevel = Level;

/**
 * Logger options
 */
export interface LoggerOptions {
  /** Include call stack in method entry logs */
  includeStack?: boolean;
  /** Maximum depth for object serialization */
  maxDepth?: number;
}

/**
 * Logger class for model debugging
 */
export class ModelLogger {
  private context: string;
  private className: string;
  private filePath?: string;
  private options: LoggerOptions;
  private entryTimestamps: Map<string, number> = new Map();

  constructor(context: string, className?: string, filePath?: string, options: LoggerOptions = {}) {
    this.context = context;
    this.className = className || context;
    this.filePath = filePath;
    this.options = {
      includeStack: false,
      maxDepth: 3,
      ...options,
    };
  }

  /**
   * Log debug information
   */
  async debug(message: string, data?: any) {
    const fullMessage = `[${this.context}] ${message}`;
    if (data !== undefined) {
      await LogService.Debug(fullMessage, 'data', JSON.stringify(data, null, 2));
    } else {
      await LogService.Debug(fullMessage);
    }
  }

  /**
   * Log info information
   */
  async info(message: string, data?: any) {
    const fullMessage = `[${this.context}] ${message}`;
    if (data !== undefined) {
      await LogService.Info(fullMessage, 'data', JSON.stringify(data, null, 2));
    } else {
      await LogService.Info(fullMessage);
    }
  }

  /**
   * Log warnings
   */
  async warn(message: string, data?: any) {
    const fullMessage = `[${this.context}] ${message}`;
    if (data !== undefined) {
      await LogService.Warning(fullMessage, 'data', JSON.stringify(data, null, 2));
    } else {
      await LogService.Warning(fullMessage);
    }
  }

  /**
   * Log errors
   */
  async error(message: string, error?: any, data?: any) {
    const fullMessage = `[${this.context}] ${message}`;
    const args: any[] = [];
    
    if (error) {
      args.push('error', error instanceof Error ? error.message : String(error));
      if (error instanceof Error && error.stack) {
        args.push('stack', error.stack);
      }
    }
    
    if (data !== undefined) {
      args.push('data', JSON.stringify(data, null, 2));
    }
    
    await LogService.Error(fullMessage, ...args);
  }

  /**
   * Log method entry with parameters
   * Includes timestamp, duration tracking, class/file info, and optional stack trace
   */
  async methodEntry(methodName: string, params?: any) {
    const timestamp = new Date().toISOString();
    const startTime = performance.now();
    
    // Store entry time for duration calculation
    this.entryTimestamps.set(methodName, startTime);
    
    // Build full method signature
    const fullMethod = `${this.className}.${methodName}`;
    
    const args: any[] = [
      'timestamp', timestamp,
      'class', this.className,
      'method', methodName,
      'fullMethod', fullMethod,
    ];
    
    // Add file path if provided
    if (this.filePath) {
      args.push('path', this.filePath);
    }
    
    // Add parameters if provided
    if (params !== undefined) {
      const paramStr = this.serializeValue(params);
      args.push('params', paramStr);
      
      // Log individual important params for easier filtering
      if (typeof params === 'object' && params !== null) {
        if (params.userId) args.push('userId', params.userId);
        if (params.id) args.push('id', params.id);
        if (params.sessionId) args.push('sessionId', params.sessionId);
        if (params.topicId) args.push('topicId', params.topicId);
      }
    }
    
    // Add stack trace if enabled
    if (this.options.includeStack) {
      const stack = new Error().stack || '';
      const callerLine = stack.split('\n')[3]?.trim() || 'unknown';
      args.push('caller', callerLine);
    }
    
    await LogService.Debug(
      `[${this.context}] → ENTER: ${fullMethod}`,
      ...args
    );
  }

  /**
   * Log method exit with result and duration
   * Includes timestamp, duration, class/file info, and return value summary
   */
  async methodExit(methodName: string, result?: any) {
    const timestamp = new Date().toISOString();
    const startTime = this.entryTimestamps.get(methodName);
    
    // Build full method signature
    const fullMethod = `${this.className}.${methodName}`;
    
    const args: any[] = [
      'timestamp', timestamp,
      'class', this.className,
      'method', methodName,
      'fullMethod', fullMethod,
    ];
    
    // Add file path if provided
    if (this.filePath) {
      args.push('path', this.filePath);
    }
    
    // Calculate duration if entry was logged
    if (startTime !== undefined) {
      const duration = performance.now() - startTime;
      args.push('duration', `${duration.toFixed(2)}ms`);
      this.entryTimestamps.delete(methodName);
    }
    
    // Add result summary if provided
    if (result !== undefined) {
      const resultStr = this.serializeValue(result);
      args.push('result', resultStr);
      
      // Log specific result fields for easier filtering
      if (typeof result === 'object' && result !== null) {
        if (result.count !== undefined) args.push('count', String(result.count));
        if (result.id) args.push('id', result.id);
        if (Array.isArray(result)) args.push('arrayLength', String(result.length));
      }
    }
    
    await LogService.Debug(
      `[${this.context}] ← EXIT: ${fullMethod}`,
      ...args
    );
  }
  
  /**
   * Serialize value for logging with depth control
   */
  private serializeValue(value: any, depth: number = 0): string {
    if (value === null) return 'null';
    if (value === undefined) return 'undefined';
    
    const maxDepth = this.options.maxDepth || 3;
    
    if (depth >= maxDepth) {
      return typeof value === 'object' ? '[Object...]' : String(value);
    }
    
    try {
      if (typeof value === 'string') return value;
      if (typeof value === 'number' || typeof value === 'boolean') return String(value);
      
      if (Array.isArray(value)) {
        if (value.length === 0) return '[]';
        if (value.length > 10) {
          return `[Array(${value.length})]`;
        }
        return JSON.stringify(value.slice(0, 10), null, 2);
      }
      
      if (typeof value === 'object') {
        const keys = Object.keys(value);
        if (keys.length === 0) return '{}';
        if (keys.length > 20) {
          return `{Object with ${keys.length} keys}`;
        }
        return JSON.stringify(value, null, 2);
      }
      
      return String(value);
    } catch (error) {
      return `[Serialization Error: ${error}]`;
    }
  }

  /**
   * Log method error
   */
  async methodError(methodName: string, error: any, params?: any) {
    await this.error(`✗ ${methodName} failed`, error, params);
  }
}

/**
 * Create a logger for a specific model
 * @param modelName - Name of the model (e.g., 'Session')
 * @param className - Class name for detailed logging (e.g., 'SessionModel')
 * @param filePath - File path relative to src (e.g., 'database/models/session')
 * @param options - Logger options
 */
export function createModelLogger(
  modelName: string, 
  className?: string, 
  filePath?: string,
  options?: LoggerOptions
): ModelLogger {
  const fullClassName = className || `${modelName}Model`;
  return new ModelLogger(`Model:${modelName}`, fullClassName, filePath, options);
}

/**
 * Create a logger for a specific service
 * @param serviceName - Name of the service
 * @param className - Class name for detailed logging
 * @param filePath - File path relative to src
 * @param options - Logger options
 */
export function createServiceLogger(
  serviceName: string, 
  className?: string,
  filePath?: string,
  options?: LoggerOptions
): ModelLogger {
  const fullClassName = className || `${serviceName}Service`;
  return new ModelLogger(`Service:${serviceName}`, fullClassName, filePath, options);
}

