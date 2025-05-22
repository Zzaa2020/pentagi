# Understanding Pentagi: How It Works with AI

Welcome to a behind-the-scenes look at Pentagi! This document explains how Pentagi is built and how it cleverly uses Artificial Intelligence (AI) to understand your requests and get things done. We'll explore its main parts, how it takes a task from start to finish, and the different 'hats' the AI wears to help out.

## Pentagi: A Look Under the Hood

Pentagi is built to be flexible and grow with new demands, much like a set of building blocks that can be rearranged and added to. It keeps different jobs separate, making it easier to manage. Think of it as having two main parts: what you see and interact with (the **Frontend**), and the engine that does the work (the **Backend**).

### Frontend: Your Control Panel

The Frontend is everything you, as a user, interact with. It's the screen where you tell Pentagi what to do, see how it's going, and get your results. When you give Pentagi a job, the Frontend sends your instructions to the Backend.

### Backend: The Engine Room

The Backend is where the magic happens. It takes your instructions and manages the AI to get the job done. Imagine it as an operations center with different managers and specialists.

Key parts of the Backend include:

*   **Router**: Think of the Router as a **receptionist**. When a request arrives from the Frontend, the Router figures out where it needs to go and directs it to the right manager.
*   **Controllers**: These are the managers overseeing different stages of a job.
    *   **FlowController (The Project Manager)**: This is like the main **project manager** for a big assignment (a "flow"). It oversees the entire job from start to finish.
    *   **TaskController (The Team Lead)**: For each major part of the job, the FlowController assigns a TaskController. This **team lead** is responsible for breaking down that part into smaller, specific steps.
    *   **SubtaskController (The Specialist)**: Each specific step is handled by a SubtaskController. This **specialist** focuses on getting one particular action done, often working directly with the AI.

### Key Supporting Components: The Backend's Toolkit

To do its work, the Backend relies on a few essential helpers:

*   **AI Provider (The Universal AI Translator)**: Pentagi might work with different AI models (like different "brains," e.g., from OpenAI or Anthropic). The AI Provider acts like a **universal translator**, allowing Pentagi to communicate clearly with whichever AI brain is chosen for the task.
*   **Prompter (The AI's Scriptwriter)**: To get the best results from an AI, you need to ask it questions or give it instructions very precisely. The Prompter is like an **AI scriptwriter** that creates these exact scripts (called "prompts") for the AI, making sure it understands what to do based on the current step, available tools, and the overall goal.
*   **ToolsExecutor (The Workshop Supervisor)**: Sometimes, the AI needs to do things in the real world, like browse a website, run a command, or check a file. The ToolsExecutor is like a **workshop supervisor**. It has a set of approved tools the AI can ask to use, and it makes sure these tools are used correctly and safely.

This setup, with clear roles and responsibilities, makes Pentagi adaptable and easier to update. Now that we've seen the main parts of Pentagi, let's look at how they work together when you give Pentagi a job.

## Lifecycle of a User-Triggered Task: From Your Request to Results

Let's walk through what happens when you give Pentagi a job, like "Scan example.com for security issues."

1.  **You Start the Job**: You type your goal into the Frontend and hit "start."

2.  **Instruction Received**: The Frontend sends your request to the Backend's **Router** (the receptionist), which sees it as a new major assignment (a "flow") and passes it to a **FlowController** (the project manager).

3.  **The Project Manager Takes Over (FlowWorker Activation)**:
    *   The **FlowController** (our Project Manager) gets a `FlowWorker` ready to handle this entire assignment.
    *   This worker sets up the necessary support:
        *   The **AI Provider** (universal translator) is chosen and configured.
        *   The **Prompter** (AI's scriptwriter) gets ready to draft instructions for the AI.
        *   The **ToolsExecutor** (workshop supervisor) prepares the set of tools the AI might need.

4.  **Breaking Down the Job (Task Creation & Decomposition)**:
    *   The `FlowWorker` (Project Manager) hands off the main goal to a `TaskWorker`, managed by a **TaskController** (our Team Lead).
    *   The `TaskWorker` (Team Lead), with help from the **Prompter**, asks the **AI Provider** to act as a *Planner*. The AI looks at your goal ("Scan example.com") and breaks it down into a list of smaller steps, like: "1. Check what services are running on example.com," "2. See if those services have known problems," etc.

5.  **Getting to Work, Step-by-Step (Subtask Execution Loop)**:
    *   The `TaskWorker` (Team Lead) takes these steps (subtasks) one by one.
    *   For each step, it brings in a `SubtaskWorker`, managed by a **SubtaskController** (our Specialist).
    *   The `SubtaskWorker` (Specialist) uses the **Prompter** to tell the **AI Provider** to act as an *Executor*. The AI now focuses on just this one step.
    *   The AI, as the *Executor*, might decide it needs a tool. For instance, to "check what services are running," it might ask the **ToolsExecutor** (workshop supervisor) for a "network scanning" tool. The results from the tool are sent back to the AI.
    *   **Need More Info?**: If the AI realizes it can't finish the step without more information from you, it tells the `SubtaskWorker`. The whole process pauses, waiting for your input.

6.  **Your Turn (Handling User Input)**:
    *   If paused, you provide the needed information through the Frontend.
    *   This info goes back through the Project Manager, Team Lead, and Specialist to the waiting AI.
    *   With the new information, the AI gets back to work on that step.

7.  **Wrapping Up a Major Part (Task Completion & Summarization)**:
    *   Once all steps under a `TaskWorker` (Team Lead) are done, it asks the **AI Provider** to act as a *Summarizer*.
    *   The AI looks at everything that happened in those steps and writes a summary for that part of the job (e.g., "Found services A, B, and C. Service A has a potential issue.").

8.  **Project Manager Reviews (Flow Continuation)**:
    *   The `FlowWorker` (Project Manager) gets this summary.
    *   Depending on how the overall assignment was set up, it might show you these results, start another major part of the job, ask for your next instruction, or decide the whole assignment is complete.

This process allows Pentagi to tackle complex goals by breaking them down, using AI in different roles, and involving you when needed. The lifecycle showed the AI playing different parts. Let's take a closer look at these specific AI roles.

## AI Roles Involved: The Different Hats AI Wears

Throughout a task, the AI, guided by the **Prompter** and communicating through the **AI Provider**, takes on several specific roles, much like different experts on a team:

1.  **Task Titler (The Quick Labeler)**:
    *   **Job**: Comes up with a short, clear title for a task based on your initial request.
    *   **When**: When a new `TaskWorker` (Team Lead) first gets an assignment. This helps everyone keep track of things.

2.  **Planner/Decomposer (The Strategist)**:
    *   **Job**: Looks at the main goal given to a `TaskWorker` (Team Lead) and figures out the best strategy by breaking it into smaller, manageable subtasks.
    *   **When**: When the `TaskWorker` (Team Lead) needs to create a step-by-step plan.

3.  **Executor/Agent (The Hands-On Expert)**:
    *   **Job**: This is the AI getting its hands dirty, focusing on one specific subtask managed by a `SubtaskWorker` (Specialist). It:
        *   Understands what needs to be done for that single step.
        *   Decides if it needs a tool from the **ToolsExecutor** (workshop supervisor) – like running a command or checking a file.
        *   Uses the tool and understands the results.
        *   Thinks about the results to decide what to do next.
        *   Figures out if it has everything it needs or if it should ask you for help.
    *   **When**: Whenever a `SubtaskWorker` (Specialist) is working on a specific action. This is where most of the "doing" happens.

4.  **Summarizer/Synthesizer (The Reporter)**:
    *   **Job**: After a `TaskWorker` (Team Lead) has overseen the completion of all its subtasks, this AI role steps in to look at all the individual results, logs, and outcomes. It then writes up a clear, overall summary of what was achieved for that entire task.
    *   **When**: When a `TaskWorker` (Team Lead) finishes all its assigned steps and needs to report back the final outcome.

By having the AI switch between these roles, Pentagi can use its intelligence effectively at every stage – from planning the work to doing it, and finally, to reporting back on how it all went. The system's "managers" (Controllers and Workers) make sure the AI gets the right script (prompt) for the right role.

---

In essence, Pentagi works like a well-organized team. It combines traditional software components—like project managers, team leads, specialists, and tool supervisors—with a versatile AI that can act as a planner, a hands-on expert, and a reporter. This teamwork allows Pentagi to tackle complex tasks in a smart and adaptable way, always keeping your goals in focus.
