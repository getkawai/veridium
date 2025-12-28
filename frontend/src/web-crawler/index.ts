// Web crawler functionality for the application

export interface CrawlSuccessResult {
  title: string;
  content: string;
  url: string;
  website: string;
}

export interface CrawlErrorResult {
  errorMessage: string;
  errorType: string;
  url: string;
}