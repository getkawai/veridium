export default function DocsAPIImage() {
  return (
    <div>
      <div className="page-header">
        <h2>Image Generation API</h2>
        <p>Generate, edit, and create variations of images. Compatible with the OpenAI Images API.</p>
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

          <div className="card" id="image-generation">
            <h3>Image Generation</h3>
            <p>Generate images from text prompts.</p>

            <div className="doc-section" id="image-generation-post--images-generations">
              <h4><span className="method-post">POST</span> /images/generations</h4>
              <p className="doc-description">Generates images from a text prompt.</p>
              <p><strong>Authentication:</strong> Required when auth is enabled. Token must have 'images-generations' endpoint access.</p>
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
                    <td><code>prompt</code></td>
                    <td><code>string</code></td>
                    <td>Yes</td>
                    <td>Text description of the image to generate</td>
                  </tr>
                  <tr>
                    <td><code>model</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Model ID to use</td>
                  </tr>
                  <tr>
                    <td><code>n</code></td>
                    <td><code>integer</code></td>
                    <td>No</td>
                    <td>Number of images to generate (1-10, default: 1)</td>
                  </tr>
                  <tr>
                    <td><code>size</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Image size (e.g., '1024x1024')</td>
                  </tr>
                  <tr>
                    <td><code>quality</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Image quality: 'standard' or 'hd'</td>
                  </tr>
                  <tr>
                    <td><code>response_format</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Response format: 'url' or 'b64_json'</td>
                  </tr>
                </tbody>
              </table>
              <h5>Response</h5>
              <p>Returns generated image URLs or base64-encoded images.</p>
              <h5>Example</h5>
              <p className="example-label"><strong>Generate a single image:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/images/generations \\
  -H "Authorization: Bearer $API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "prompt": "A serene mountain landscape at sunset",
    "n": 1,
    "size": "1024x1024"
  }'`}</code>
              </pre>
            </div>
          </div>

          <div className="card" id="image-editing">
            <h3>Image Editing</h3>
            <p>Edit existing images using masks and prompts.</p>

            <div className="doc-section" id="image-editing-post--images-edits">
              <h4><span className="method-post">POST</span> /images/edits</h4>
              <p className="doc-description">Edits an image based on a text prompt and optional mask.</p>
              <p><strong>Authentication:</strong> Required when auth is enabled. Token must have 'images-edits' endpoint access.</p>
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
                    <td><code>prompt</code></td>
                    <td><code>string</code></td>
                    <td>Yes</td>
                    <td>Text description of the edit</td>
                  </tr>
                  <tr>
                    <td><code>image</code></td>
                    <td><code>string</code></td>
                    <td>Yes</td>
                    <td>Base64-encoded image to edit</td>
                  </tr>
                  <tr>
                    <td><code>mask</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Base64-encoded mask (white=edit, black=keep)</td>
                  </tr>
                  <tr>
                    <td><code>n</code></td>
                    <td><code>integer</code></td>
                    <td>No</td>
                    <td>Number of images (1-10, default: 1)</td>
                  </tr>
                  <tr>
                    <td><code>size</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Image size</td>
                  </tr>
                  <tr>
                    <td><code>response_format</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Response format: 'url' or 'b64_json'</td>
                  </tr>
                </tbody>
              </table>
              <h5>Response</h5>
              <p>Returns edited image URLs or base64-encoded images.</p>
              <h5>Example</h5>
              <p className="example-label"><strong>Edit image with mask:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/images/edits \\
  -H "Authorization: Bearer $API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "prompt": "Add a beautiful sunset sky",
    "image": "data:image/png;base64,iVBORw0KG...",
    "mask": "data:image/png;base64,iVBORw0KG...",
    "n": 1
  }'`}</code>
              </pre>
            </div>
          </div>

          <div className="card" id="image-variations">
            <h3>Image Variations</h3>
            <p>Create variations of existing images.</p>

            <div className="doc-section" id="image-variations-post--images-variations">
              <h4><span className="method-post">POST</span> /images/variations</h4>
              <p className="doc-description">Creates variations of an input image.</p>
              <p><strong>Authentication:</strong> Required when auth is enabled. Token must have 'images-variations' endpoint access.</p>
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
                    <td><code>image</code></td>
                    <td><code>string</code></td>
                    <td>Yes</td>
                    <td>Base64-encoded image</td>
                  </tr>
                  <tr>
                    <td><code>n</code></td>
                    <td><code>integer</code></td>
                    <td>No</td>
                    <td>Number of variations (1-10)</td>
                  </tr>
                  <tr>
                    <td><code>size</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Image size</td>
                  </tr>
                  <tr>
                    <td><code>response_format</code></td>
                    <td><code>string</code></td>
                    <td>No</td>
                    <td>Response format: 'url' or 'b64_json'</td>
                  </tr>
                </tbody>
              </table>
              <h5>Response</h5>
              <p>Returns variation image URLs or base64-encoded images.</p>
              <h5>Example</h5>
              <p className="example-label"><strong>Generate 4 variations:</strong></p>
              <pre className="code-block">
                <code>{`curl -X POST https://api.getkawai.com/v1/images/variations \\
  -H "Authorization: Bearer $API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "image": "data:image/png;base64,iVBORw0KG...",
    "n": 4
  }'`}</code>
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
              <a href="#image-generation" className="doc-index-header">Image Generation</a>
              <ul>
                <li><a href="#image-generation-post--images-generations">POST /images/generations</a></li>
              </ul>
            </div>
            <div className="doc-index-section">
              <a href="#image-editing" className="doc-index-header">Image Editing</a>
              <ul>
                <li><a href="#image-editing-post--images-edits">POST /images/edits</a></li>
              </ul>
            </div>
            <div className="doc-index-section">
              <a href="#image-variations" className="doc-index-header">Image Variations</a>
              <ul>
                <li><a href="#image-variations-post--images-variations">POST /images/variations</a></li>
              </ul>
            </div>
          </div>
        </nav>
      </div>
    </div>
  );
}
