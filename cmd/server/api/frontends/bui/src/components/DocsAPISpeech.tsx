export default function DocsAPISpeech() {
  return (
    <div>
      <div className="page-header">
        <h2>Text-to-Speech API</h2>
        <p>Generate speech from text using TTS models. Compatible with the OpenAI Speech API.</p>
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

          <div className="card" id="text-to-speech">
            <h3>Text-to-Speech</h3>
            <p>Generate audio from text using text-to-speech models.</p>

            <div className="doc-section" id="text-to-speech-post--audio-speech">
              <h4><span className="method-post">POST</span> /audio/speech</h4>
              <p className="doc-description">Generates audio from the input text using a specified voice. Supports multiple output formats.</p>
              <p><strong>Authentication:</strong> Required when auth is enabled. Token must have 'audio-speech' endpoint access.</p>
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
                    <td>Must be application/json</td>
                  </tr>
                </tbody>
              </table>
              <h5>Request Body</h5>
              <p><code>application/json</code></p>
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
                    <td><code>model</code></td>
                    <td><code>string</code></td>
                    <td>Yes</td>
                    <td>TTS model ID (e.g., 'kokoro')</td>
                  </tr>
                  <tr>
                    <td><code>input</code></td>
                    <td><code>string</code></td>
                    <td>Yes</td>
                    <td>The text to generate audio for</td>
                  </tr>
                  <tr>
                    <td><code>voice</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>The voice to use (default: 'af_sarah' for Kokoro)</td>
                  </tr>
                  <tr>
                    <td><code>response_format</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Audio format: mp3, opus, aac, flac, wav, pcm (default: mp3)</td>
                  </tr>
                  <tr>
                    <td><code>speed</code></td>
                    <td><code>number</code></td>
                    <td>No</td>
                    <td>Speech speed multiplier 0.25-4.0 (default: 1.0)</td>
                  </tr>
                </tbody>
              </table>
              <h5>Response</h5>
              <p>Returns binary audio data in the requested format.</p>
              <h5>Example</h5>
              <p className="example-label"><strong>Generate speech with default voice:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/audio/speech \\
  -H "Authorization: Bearer $API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "kokoro",
    "input": "Hello, this is a test of text to speech.",
    "voice": "af_sarah"
  }' \\
  --output speech.mp3`}</code>
              </pre>
              <p className="example-label"><strong>Generate speech with specific voice and format:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/audio/speech \\
  -H "Authorization: Bearer $API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "kokoro",
    "input": "Welcome to Kawai DeAI Network!",
    "voice": "af_sarah",
    "response_format": "wav"
  }' \\
  --output speech.wav`}</code>
              </pre>
              <p className="example-label"><strong>Generate speech with slower speed:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/audio/speech \\
  -H "Authorization: Bearer $API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "kokoro",
    "input": "Please listen carefully to these instructions.",
    "voice": "af_sarah",
    "speed": 0.8
  }' \\
  --output speech.mp3`}</code>
              </pre>
            </div>
          </div>

          <div className="card" id="supported-models">
            <h3>Supported Models</h3>
            <p>Available TTS models and voices for speech generation.</p>

            <div className="doc-section" id="supported-models--kokoro-tts">
              <h4>Kokoro TTS</h4>
              <p className="doc-description">Kokoro is an open-source TTS model optimized for high-quality English speech synthesis. It supports multiple voices.</p>
              <h5>Example</h5>
              <p className="example-label"><strong>Available voices for Kokoro:</strong></p>
              <pre className="code-block">
                <code>{`// American English voices
af_sarah    - Sarah (Female)
af_nicole   - Nicole (Female)
am_adam     - Adam (Male)

// British English voices
bf_emma     - Emma (Female)
bm_george   - George (Male)

// Other English voices
af_bella    - Bella (Female)
af_heart    - Heart (Female)`}</code>
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
              <a href="#text-to-speech" className="doc-index-header">Text-to-Speech</a>
              <ul>
                <li><a href="#text-to-speech-post--audio-speech">POST /audio/speech</a></li>
              </ul>
            </div>
            <div className="doc-index-section">
              <a href="#supported-models" className="doc-index-header">Supported Models</a>
              <ul>
                <li><a href="#supported-models--kokoro-tts">Kokoro TTS</a></li>
              </ul>
            </div>
          </div>
        </nav>
      </div>
    </div>
  );
}
