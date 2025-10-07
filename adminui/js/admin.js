// Admin UI JavaScript for array manipulation

document.addEventListener('DOMContentLoaded', function() {
    console.log('Admin UI loaded');
});

// Add a new item to an array
function addArrayItem(arrayName) {
    const arraySection = document.querySelector(`.array-section[data-array-name="${arrayName}"]`);
    if (!arraySection) {
        console.error(`Array section for ${arrayName} not found`);
        return;
    }
    
    // Count existing items to determine the new index
    const existingItems = arraySection.querySelectorAll('.array-item');
    const newIndex = existingItems.length;
    
    // Create new item element
    const newItem = document.createElement('div');
    newItem.className = 'array-item';
    newItem.setAttribute('data-index', newIndex);
    
    // Create controls for the new item
    const controls = document.createElement('div');
    controls.className = 'array-item-controls';
    controls.innerHTML = `
        <button type="button" class="array-move-up-btn" onclick="moveArrayItemUp('${arrayName}', ${newIndex})">&#8593;</button>
        <button type="button" class="array-move-down-btn" onclick="moveArrayItemDown('${arrayName}', ${newIndex})">&#8595;</button>
        <button type="button" class="array-remove-btn" onclick="removeArrayItem('${arrayName}', ${newIndex})">Remove</button>
    `;
    newItem.appendChild(controls);
    
    // Create input field for the new item
    const inputField = document.createElement('div');
    inputField.className = 'form-field';
    inputField.innerHTML = `<label>New Item:</label><input type="text" name="${arrayName}[${newIndex}]" value="" />`;
    newItem.appendChild(inputField);
    
    // Append the new item to the array section
    arraySection.appendChild(newItem);
    
    // Update indices of all items
    updateArrayIndices(arrayName);
}

// Move an array item up
function moveArrayItemUp(arrayName, index) {
    if (index <= 0) return; // Can't move the first item up
    
    const arraySection = document.querySelector(`.array-section[data-array-name="${arrayName}"]`);
    if (!arraySection) {
        console.error(`Array section for ${arrayName} not found`);
        return;
    }
    
    const items = Array.from(arraySection.querySelectorAll('.array-item'));
    if (index >= items.length) return;
    
    const currentItem = items[index];
    const previousItem = items[index - 1];
    
    // Swap positions in the DOM
    arraySection.insertBefore(currentItem, previousItem);
    
    // Update indices of all items
    updateArrayIndices(arrayName);
}

// Move an array item down
function moveArrayItemDown(arrayName, index) {
    const arraySection = document.querySelector(`.array-section[data-array-name="${arrayName}"]`);
    if (!arraySection) {
        console.error(`Array section for ${arrayName} not found`);
        return;
    }
    
    const items = Array.from(arraySection.querySelectorAll('.array-item'));
    if (index >= items.length - 1) return; // Can't move the last item down
    
    const currentItem = items[index];
    const nextItem = items[index + 1];
    
    // Swap positions in the DOM
    arraySection.insertBefore(nextItem, currentItem);
    
    // Update indices of all items
    updateArrayIndices(arrayName);
}

// Remove an array item
function removeArrayItem(arrayName, index) {
    const arraySection = document.querySelector(`.array-section[data-array-name="${arrayName}"]`);
    if (!arraySection) {
        console.error(`Array section for ${arrayName} not found`);
        return;
    }
    
    const items = Array.from(arraySection.querySelectorAll('.array-item'));
    if (index >= items.length) return;
    
    // Remove the item
    items[index].remove();
    
    // Update indices of remaining items
    updateArrayIndices(arrayName);
}

// Update indices of all array items
function updateArrayIndices(arrayName) {
    const arraySection = document.querySelector(`.array-section[data-array-name="${arrayName}"]`);
    if (!arraySection) {
        console.error(`Array section for ${arrayName} not found`);
        return;
    }
    
    const items = Array.from(arraySection.querySelectorAll('.array-item'));
    items.forEach((item, i) => {
        item.setAttribute('data-index', i);
        
        // Update button onclick handlers
        const moveUpBtn = item.querySelector('.array-move-up-btn');
        const moveDownBtn = item.querySelector('.array-move-down-btn');
        const removeBtn = item.querySelector('.array-remove-btn');
        
        if (moveUpBtn) {
            moveUpBtn.onclick = () => moveArrayItemUp(arrayName, i);
        }
        
        if (moveDownBtn) {
            moveDownBtn.onclick = () => moveArrayItemDown(arrayName, i);
        }
        
        if (removeBtn) {
            removeBtn.onclick = () => removeArrayItem(arrayName, i);
        }
        
        // Update input names
        const inputs = item.querySelectorAll('input[name]');
        inputs.forEach(input => {
            const name = input.getAttribute('name');
            // Replace the index in the input name
            const newName = name.replace(/\[\d+\]/, `[${i}]`);
            input.setAttribute('name', newName);
            
            // Update labels if they exist
            const label = item.querySelector('label');
            if (label) {
                label.textContent = label.textContent.replace(/Item \d+|New Item/, `Item ${i}`);
            }
        });
    });
}