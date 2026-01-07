# Using AI Chat

Master the AI chat feature and start earning while you learn, code, and create!

## 🎯 What is AI Chat?

Kawai's AI Chat gives you access to powerful language models that can help you with:

- 💻 **Code**: Write, debug, and explain code
- ✍️ **Writing**: Create content, emails, documents
- 🎓 **Learning**: Explain concepts, answer questions
- 🔍 **Research**: Summarize information, analyze data
- 🎨 **Creative**: Brainstorm ideas, write stories

**Best part?** You earn 5% cashback on every request!

## 🚀 Getting Started

### Opening the Chat

1. Launch Kawai Desktop App
2. Main interface = Chat (default view)
3. Click **"New Chat"** to start fresh

### Your First Message

1. Type your question or request
2. Press **Enter** or click **Send**
3. Wait for AI response (usually 2-10 seconds)
4. Continue the conversation!

**Example:**
```
You: What is blockchain?
AI: Blockchain is a distributed ledger technology...
```

## 💬 Chat Features

### Multi-Turn Conversations

The AI remembers your conversation context.

**Example:**
```
You: Write a function to calculate Fibonacci numbers
AI: [provides Python code]

You: Can you make it more efficient?
AI: [provides optimized version with memoization]

You: Now convert it to JavaScript
AI: [provides JS version]
```

**Tip:** Each conversation has full context - no need to repeat yourself!

### Code Highlighting

Code in responses is automatically highlighted with syntax.

**Supported languages:**
- Python, JavaScript, TypeScript
- Go, Rust, C++, Java
- HTML, CSS, SQL
- And 50+ more!

### Markdown Support

AI responses support full Markdown formatting:

- **Bold** and *italic* text
- Lists (bullet and numbered)
- Tables
- Inline `code` and code blocks
- Links and images
- Headers and sections

### Copy & Export

- **Copy code**: Click the copy icon on code blocks
- **Copy message**: Right-click any message
- **Export chat**: Save conversations to file (coming soon)

## 🎨 Chat Management

### Creating New Chats

Click **"New Chat"** to start a fresh conversation.

**When to start new:**
- Different topic/project
- Conversation getting too long
- Want fresh context

### Chat History

All your chats are saved automatically.

**Accessing history:**
1. Click the **history icon** (sidebar)
2. Browse your past conversations
3. Click any chat to continue

**Features:**
- ✅ Auto-saved (no manual save needed)
- ✅ Searchable (find old conversations)
- ✅ Organized by date
- ✅ Unlimited storage

### Renaming Chats

1. Click the **pencil icon** next to chat title
2. Enter new name
3. Press Enter

**Tip:** Use descriptive names like "Marketing Campaign Ideas" or "Python Tutorial".

### Deleting Chats

1. Right-click the chat
2. Select **"Delete"**
3. Confirm

**Note:** Deleted chats are gone forever. Be sure!

## 🧠 Model Selection

### Available Models

Currently using local LLM models via `llama.cpp`:

- **Fast Mode**: Quick responses, lower quality
- **Balanced**: Good speed and quality (default)
- **Quality Mode**: Best responses, slower

### Switching Models

1. Click **settings icon** in chat
2. Select **"Model"**
3. Choose your preferred model

**When to use each:**
- **Fast**: Quick questions, simple tasks
- **Balanced**: Most use cases
- **Quality**: Complex reasoning, long content

## 📚 Knowledge Base Integration

### What is Knowledge Base?

Upload documents so AI can reference your files.

**Use cases:**
- Code repositories → AI helps with your codebase
- Documentation → AI answers from your docs
- Research papers → AI summarizes your research
- Business reports → AI analyzes your data

### Adding Documents

1. Click **"Knowledge Base"** in sidebar
2. Click **"Upload"**
3. Select files:
   - PDF documents
   - Word files (.docx)
   - Markdown (.md)
   - Text files (.txt)
   - Code files

4. Wait for indexing (automatic)
5. Ask questions about your docs!

### Using Knowledge Base

Once documents are uploaded, just ask questions:

**Example:**
```
You: What were the Q3 sales numbers?
AI: According to the Q3_Report.pdf you uploaded,
     sales were $2.5M, up 25% from Q2...
```

**Pro tip:** AI automatically knows which documents are relevant!

### Managing Documents

- **View all**: Click "Knowledge Base" to see all docs
- **Remove**: Click X on any document
- **Update**: Upload newer version (same filename)

## 💰 Earning While You Chat

### Use-to-Earn (5% Cashback)

Every AI request earns you KAWAI tokens!

**Example:**
```
Request cost: 0.05 USDT
Cashback: 0.05 × 5% = 2.5 KAWAI tokens
```

**The more you use, the more you earn!**

### Tracking Your Earnings

1. Go to **Wallet → Rewards → Use-to-Earn**
2. See your accumulated rewards
3. Claim weekly after settlement

### Cost Breakdown

**Average costs:**
- Simple question: 0.01-0.02 USDT
- Code generation: 0.03-0.05 USDT
- Long analysis: 0.05-0.10 USDT

**With 100 USDT balance:**
- ~2,000-10,000 AI requests
- Earn 250-500 KAWAI tokens back
- Net cost: 95 USDT (after cashback)

## 🎯 Pro Tips

### Getting Better Responses

**1. Be specific**
```
❌ "Write code"
✅ "Write a Python function that calculates compound interest"
```

**2. Provide context**
```
❌ "Fix this"
✅ "Fix this JavaScript function - it should validate email addresses but returns false for valid emails"
```

**3. Use examples**
```
✅ "Create a function like this: input: [1,2,3], output: 6"
```

**4. Break down complex tasks**
```
✅ "First, explain the algorithm. Then, show the code."
```

### Using System Messages

System messages guide the AI's behavior.

**Examples:**
```
System: You are a Python expert. Be concise.
System: Explain to a beginner. Use simple terms.
System: Act as a code reviewer. Be critical.
```

**How to use:**
1. Click **settings** → **System Message**
2. Enter your instructions
3. Applies to entire chat

### Iterating on Responses

Don't hesitate to refine!

**Example flow:**
```
You: Write a blog post about AI
AI: [provides 500 words]

You: Make it shorter and more casual
AI: [provides 300 words, casual tone]

You: Add 3 examples
AI: [adds examples]
```

## 🔧 Advanced Features

### Code Execution (Coming Soon)

Run code directly in the chat:

```python
print("Hello, Kawai!")
# Output: Hello, Kawai!
```

### Image Understanding (Coming Soon)

Upload images and ask questions:

```
You: [uploads diagram] Explain this architecture
AI: This diagram shows a microservices architecture...
```

### Voice Input (Coming Soon)

Speak your questions instead of typing.

### Collaborative Chats (Coming Soon)

Share chats with team members.

## 📊 Performance & Limits

### Speed

- **Local inference**: 20-50 tokens/second
- **Average response time**: 2-10 seconds
- **Long responses**: up to 30 seconds

**Factors affecting speed:**
- Your hardware (CPU/GPU)
- Model size
- Response length
- System load

### Token Limits

- **Input**: ~4,000 tokens per message
- **Output**: ~2,000 tokens per response
- **Context window**: ~8,000 tokens total

**What's a token?** ~4 characters or 0.75 words.

### Rate Limits

- **No hard limits** on chat requests
- Limited only by your USDT balance
- Fair usage policy applies

## 🆘 Troubleshooting

### AI Not Responding

**Solutions:**
1. Check internet connection
2. Verify USDT balance > 0
3. Restart the app
4. Check logs: Settings → Advanced → View Logs

### Slow Responses

**Solutions:**
1. Switch to **Fast Mode**
2. Ask shorter questions
3. Close other apps (free up RAM)
4. Check system resources

### Unexpected Responses

**Solutions:**
1. Rephrase your question
2. Add more context
3. Try a different model
4. Start a new chat (fresh context)

### "Insufficient Balance"

**Solution:**
1. Check your USDT balance
2. [Deposit more USDT](deposit.md)
3. Or [claim free trial](free-trial.md)

### Context Confusion

If AI seems confused:
1. Start a **new chat**
2. Or explicitly reset: "Forget the above, new topic..."
3. Be more explicit with context

## 💡 Use Cases & Examples

### For Developers

**Code review:**
```
You: Review this code: [paste code]
AI: [provides detailed review with suggestions]
```

**Debugging:**
```
You: Why is this throwing a null pointer exception?
AI: [analyzes and explains the issue]
```

**Documentation:**
```
You: Write documentation for this function
AI: [generates comprehensive docs]
```

### For Content Creators

**Blog posts:**
```
You: Write a 500-word blog about Web3
AI: [provides article]
```

**Social media:**
```
You: Create 5 tweets about our new product
AI: [provides tweet thread]
```

**Video scripts:**
```
You: Script for a 2-minute YouTube intro
AI: [provides script with timing]
```

### For Students

**Explanations:**
```
You: Explain quantum computing like I'm 10
AI: [simple explanation]
```

**Study help:**
```
You: Quiz me on photosynthesis
AI: [asks questions]
```

**Essay help:**
```
You: Outline for essay on climate change
AI: [provides structured outline]
```

### For Business

**Emails:**
```
You: Write a professional follow-up email
AI: [drafts email]
```

**Analysis:**
```
You: Analyze this sales data [paste data]
AI: [provides insights]
```

**Proposals:**
```
You: Create project proposal outline
AI: [provides structure]
```

## ✅ Checklist

For best results:

- [ ] Have USDT balance ready
- [ ] Understand token costs
- [ ] Know how to use Knowledge Base
- [ ] Practice with clear, specific prompts
- [ ] Track your cashback rewards
- [ ] Save important chats

## 🚀 Next Steps

Explore more features:

1. **[Image Generation](image-generation.md)** - Create AI art
2. **[Rewards Dashboard](../rewards/overview.md)** - Track earnings
3. **[Referral Program](../rewards/referral.md)** - Earn by sharing
4. **[Trading](../trading/marketplace.md)** - Trade KAWAI tokens

---

**Need help?** Check [FAQ](../faq/general.md) or [join Discord](https://discord.gg/kawai).

