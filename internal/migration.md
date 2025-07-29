
## Repository Migration Guide

### Moving to BINAR-Learning Organization

1. Clone the demo repository:
   ```bash
   # Clone demo repository
   git clone https://github.com/BINAR-Learning/demo-repository.git
   cd demo-repository
   ```

2. Create and switch to new branch:
   ```bash
   # Create new branch for your work
   git checkout -b feature/add-stock-consolidation
   ```

3. Create your project directory:
   ```bash
   # For Windows PowerShell/CMD
   mkdir ".\[INDIVIDU 12] - [Stock Consolidation]"
   cd ".\[INDIVIDU 12] - [Stock Consolidation]"

   # Alternative command if above doesn't work
   cd "[INDIVIDU 12] - [Stock Consolidation]"
   # OR
   cd [INDIVIDU" "12]" "-" "[Stock" "Consolidation]
   ```

4. Clone your stock consolidation project:
   ```bash
   # Clone your project into the new directory
   git clone https://github.com/talahoo/stock-consolidation.git
   cd stock-consolidation
   # Remove git history (optional)
   rm -rf .git
   ```

5. Commit and push changes:
   ```bash
   # Go back to demo-repository root
   cd ../..
   # Add changes
   git add "[INDIVIDU 12] - [Stock Consolidation]"
   git commit -m "feat: add stock consolidation project"
   # Push to remote
   git push origin feature/add-stock-consolidation
   ```

6. Create Pull Request:
   - Go to https://github.com/BINAR-Learning/demo-repository
   - Click "Compare & pull request" for your branch
   - Set title: "feat: [INDIVIDU 12] - [Stock Consolidation]"
   - Add description about your project
   - Create pull request

7. After PR is merged:
   - You can archive or delete your original repository
   - Update any documentation references to point to new location

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
