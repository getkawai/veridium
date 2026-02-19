export default function DocsAPITranscriptions() {
  return (
    <div>
      <div className="page-header">
        <h2>Speech-to-Text API</h2>
        <p>Transcribe audio to text or translate audio to English. Compatible with the OpenAI Audio API.</p>
      </div>

      <div className="doc-layout">
        <div className="doc-content">
          <div className="card" id="overview">
            <h3>Overview</h3>
            <p>All endpoints are prefixed with <code>/v1</code>. Base URL: <code>https://api.getkawai.com</code></p>
            <h4>Authentication</h4>
            <p>When authentication is enabled, include the token in the Authorization header:</p>
            <pre className="code-block">
              <code>Authorization: Bearer API_KEY</code>
            </pre>
          </div>

          <div className="card" id="transcriptions">
            <h3>Transcriptions</h3>
            <p>Transcribe audio to text in the original language.</p>

            <div className="doc-section" id="transcriptions-post--audio-transcriptions">
              <h4><span className="method-post">POST</span> /audio/transcriptions</h4>
              <p className="doc-description">Transcribes audio into the input language. Supports multiple response formats including verbose JSON with timestamps.</p>
              <p><strong>Authentication:</strong> Required when auth is enabled. Token must have 'audio-transcriptions' endpoint access.</p>
              <h5>Headers</h5>
              <table className="flags-table">
                <thead>
                  <tr>
                    <th>Header</th>
                    <th>Required</th>
                    <th>Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td><code>Authorization</code></td>
                    <td>Yes</td>
                    <td>Bearer token for authentication</td>
                  </tr>
                  <tr>
                    <td><code>Content-Type</code></td>
                    <td>Yes</td>
                    <td>Must be multipart/form-data</td>
                  </tr>
                </tbody>
              </table>
              <h5>Request Body</h5>
              <p><code>multipart/form-data</code></p>
              <table className="flags-table">
                <thead>
                  <tr>
                    <th>Field</th>
                    <th>Type</th>
                    <th>Required</th>
                    <th>Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td><code>file</code></td>
                    <td><code>binary</code></td>
                    <td>Yes</td>
                    <td>Audio file (flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, webm)</td>
                  </tr>
                  <tr>
                    <td><code>model</code></td>
                    <td><code>string</code></td>
                    <td>Yes</td>
                    <td>Transcription model (e.g., 'tiny', 'base', 'small', 'medium', 'large')</td>
                  </tr>
                  <tr>
                    <td><code>language</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Language code (ISO-639-1). Auto-detected if not provided.</td>
                  </tr>
                  <tr>
                    <td><code>prompt</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Optional text to guide style or continue previous segment</td>
                  </tr>
                  <tr>
                    <td><code>response_format</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Format: json, text, srt, vtt, verbose_json (default: json)</td>
                  </tr>
                  <tr>
                    <td><code>temperature</code></td>
                    <td><code>number</code></td>
                    <td>No</td>
                    <td>Sampling temperature 0-1 (default: 0)</td>
                  </tr>
                </tbody>
              </table>
              <h5>Response</h5>
              <p>Returns transcription text. Verbose JSON includes segments, timestamps, and language detection.</p>
              <h5>Example</h5>
              <p className="example-label"><strong>Basic transcription:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/audio/transcriptions \\
  -H "Authorization: Bearer $API_KEY" \\
  -F "file=@audio.mp3" \\
  -F "model=base" \\
  -F "language=en"`}</code>
              </pre>
              <p className="example-label"><strong>Verbose JSON with timestamps:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/audio/transcriptions \\
  -H "Authorization: Bearer $API_KEY" \\
  -F "file=@audio.mp3" \\
  -F "model=base" \\
  -F "response_format=verbose_json"`}</code>
              </pre>
            </div>
          </div>

          <div className="card" id="translations">
            <h3>Translations</h3>
            <p>Translate audio from any language to English.</p>

            <div className="doc-section" id="translations-post--audio-translations">
              <h4><span className="method-post">POST</span> /audio/translations</h4>
              <p className="doc-description">Translates audio into English. The source language is automatically detected.</p>
              <p><strong>Authentication:</strong> Required when auth is enabled. Token must have 'audio-translations' endpoint access.</p>
              <h5>Headers</h5>
              <table className="flags-table">
                <thead>
                  <tr>
                    <th>Header</th>
                    <th>Required</th>
                    <th>Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td><code>Authorization</code></td>
                    <td>Yes</td>
                    <td>Bearer token for authentication</td>
                  </tr>
                  <tr>
                    <td><code>Content-Type</code></td>
                    <td>Yes</td>
                    <td>Must be multipart/form-data</td>
                  </tr>
                </tbody>
              </table>
              <h5>Request Body</h5>
              <p><code>multipart/form-data</code></p>
              <table className="flags-table">
                <thead>
                  <tr>
                    <th>Field</th>
                    <th>Type</th>
                    <th>Required</th>
                    <th>Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td><code>file</code></td>
                    <td><code>binary</code></td>
                    <td>Yes</td>
                    <td>Audio file (flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, webm)</td>
                  </tr>
                  <tr>
                    <td><code>model</code></td>
                    <td><code>string</code></td>
                    <td>Yes</td>
                    <td>Translation model (e.g., 'tiny', 'base', 'small', 'medium', 'large')</td>
                  </tr>
                  <tr>
                    <td><code>prompt</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Optional text to guide style</td>
                  </tr>
                  <tr>
                    <td><code>response_format</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Format: json, text, srt, vtt, verbose_json (default: json)</td>
                  </tr>
                  <tr>
                    <td><code>temperature</code></td>
                    <td><code>number</code></td>
                    <td>No</td>
                    <td>Sampling temperature 0-1 (default: 0)</td>
                  </tr>
                </tbody>
              </table>
              <h5>Response</h5>
              <p>Returns English translation text. Verbose JSON includes segments and timestamps.</p>
              <h5>Example</h5>
              <p className="example-label"><strong>Translate Spanish audio to English:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/audio/translations \\
  -H "Authorization: Bearer $API_KEY" \\
  -F "file=@spanish-audio.mp3" \\
  -F "model=base"`}</code>
              </pre>
              <p className="example-label"><strong>Translate with verbose output:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/audio/translations \\
  -H "Authorization: Bearer $API_KEY" \\
  -F "file=@french-audio.mp3" \\
  -F "model=base" \\
  -F "response_format=verbose_json"`}</code>
              </pre>
            </div>
          </div>

          <div className="card" id="response-formats">
            <h3>Response Formats</h3>
            <p>Available response formats for transcription and translation.</p>

            <div className="doc-section" id="response-formats--json-(default)">
              <h4>JSON (default)</h4>
              <p className="doc-description">Simple JSON response with only the transcribed/translated text.</p>
              <h5>Example</h5>
              <pre className="code-block">
                <code>{`{
  "text": "Hello, this is the transcribed text."
}`}</code>
              </pre>
            </div>

            <div className="doc-section" id="response-formats--verbose-json">
              <h4>Verbose JSON</h4>
              <p className="doc-description">Detailed JSON response with language detection, duration, segments, and word-level timestamps.</p>
              <h5>Example</h5>
              <pre className="code-block">
                <code>{`{
  "task": "transcribe",
  "language": "en",
  "duration": 5.2,
  "text": "Hello, this is the transcribed text.",
  "segments": [
    {
      "id": 0,
      "start": 0.0,
      "end": 2.5,
      "text": "Hello, this is",
      "tokens": [123, 456, 789]
    }
  ],
  "words": [
    {"word": "Hello", "start": 0.0, "end": 0.5}
  ]
}`}</code>
              </pre>
            </div>
          </div>

          <div className="card" id="supported-models">
            <h3>Supported Models</h3>
            <p>Whisper models available for transcription and translation.</p>

            <div className="doc-section" id="supported-models--whisper-models">
              <h4>Whisper Models</h4>
              <p className="doc-description">OpenAI Whisper models for speech recognition.</p>
              <h5>Example</h5>
              <p className="example-label"><strong>Available models:</strong></p>
              <pre className="code-block">
                <code>{`tiny    - 39M parameters, fastest
base    - 74M parameters, good balance
small   - 244M parameters, better accuracy
medium  - 769M parameters, high accuracy
large   - 1550M parameters, best accuracy`}</code>
              </pre>
            </div>
          </div>
        </div>

        <nav className="doc-sidebar">
          <div className="doc-sidebar-content">
            <div className="doc-index-section">
              <a href="#overview" className="doc-index-header">Overview</a>
            </div>
            <div className="doc-index-section">
              <a href="#transcriptions" className="doc-index-header">Transcriptions</a>
              <ul>
                <li><a href="#transcriptions-post--audio-transcriptions">POST /audio/transcriptions</a></li>
              </ul>
            </div>
            <div className="doc-index-section">
              <a href="#translations" className="doc-index-header">Translations</a>
              <ul>
                <li><a href="#translations-post--audio-translations">POST /audio/translations</a></li>
              </ul>
            </div>
            <div className="doc-index-section">
              <a href="#response-formats" className="doc-index-header">Response Formats</a>
              <ul>
                <li><a href="#response-formats--json-(default)">JSON (default)</a></li>
                <li><a href="#response-formats--verbose-json">Verbose JSON</a></li>
              </ul>
            </div>
            <div className="doc-index-section">
              <a href="#supported-models" className="doc-index-header">Supported Models</a>
              <ul>
                <li><a href="#supported-models--whisper-models">Whisper Models</a></li>
              </ul>
            </div>
          </div>
        </nav>
      </div>
    </div>
  );
}
