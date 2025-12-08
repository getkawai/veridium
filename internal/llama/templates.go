/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package llama

// Llama32ToolTemplate is a custom chat template for Llama 3.2 models with tool support.
// The original Llama 3.2 template has raise_exception for multiple tool calls,
// but gonja doesn't support that function. This template provides tool call support
// without the restrictive validation, following the same structure as Qwen.
const Llama32ToolTemplate = `{%- if tools %}
{{- '<|begin_of_text|><|start_header_id|>system<|end_header_id|>

' }}
{%- if messages[0]['role'] == 'system' %}
{{- messages[0]['content'] }}
{%- else %}
{{- 'You are a helpful assistant with access to the following functions.' }}
{%- endif %}
{{- '

# Tools

You may call one or more functions to assist with the user query.

You are provided with function signatures within <tools></tools> XML tags:
<tools>' }}
{%- for tool in tools %}
{{- '
' }}
{{- tool | tojson }}
{%- endfor %}
{{- '
</tools>

For each function call, return a json object with function name and arguments within <tool_call></tool_call> XML tags:
<tool_call>
{"name": <function-name>, "arguments": <args-json-object>}
</tool_call><|eot_id|>
' }}
{%- else %}
{%- if messages[0]['role'] == 'system' %}
{{- '<|begin_of_text|><|start_header_id|>system<|end_header_id|>

' + messages[0]['content'] + '<|eot_id|>
' }}
{%- else %}
{{- '<|begin_of_text|><|start_header_id|>system<|end_header_id|>

You are a helpful assistant.<|eot_id|>
' }}
{%- endif %}
{%- endif %}
{%- for message in messages %}
{%- if (message.role == "user") or (message.role == "system" and not loop.first) or (message.role == "assistant" and not message.tool_calls) %}
{{- '<|start_header_id|>' + message.role + '<|end_header_id|>

' + message.content + '<|eot_id|>
' }}
{%- elif message.role == "assistant" %}
{{- '<|start_header_id|>' + message.role + '<|end_header_id|>
' }}
{%- if message.content %}
{{- message.content }}
{%- endif %}
{%- for tool_call in message.tool_calls %}
{%- if tool_call.function is defined %}
{%- set tool_call = tool_call.function %}
{%- endif %}
{{- '
<tool_call>
{"name": "' }}
{{- tool_call.name }}
{{- '", "arguments": ' }}
{{- tool_call.arguments | tojson }}
{{- '}
</tool_call>' }}
{%- endfor %}
{{- '<|eot_id|>
' }}
{%- elif message.role == "tool" %}
{%- if (loop.index0 == 0) or (messages[loop.index0 - 1].role != "tool") %}
{{- '<|start_header_id|>ipython<|end_header_id|>
' }}
{%- endif %}
{{- '
<tool_response>
' }}
{{- message.content }}
{{- '
</tool_response>' }}
{%- if loop.last or (messages[loop.index0 + 1].role != "tool") %}
{{- '<|eot_id|>
' }}
{%- endif %}
{%- endif %}
{%- endfor %}
{%- if add_generation_prompt %}
{{- '<|start_header_id|>assistant<|end_header_id|>

' }}
{%- endif %}`
