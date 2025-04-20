// DOM Elements
const dropZone = document.getElementById('drop-zone');
const fileInput = document.getElementById('file-input');
const uploadProgress = document.getElementById('upload-progress');
const progressBarFill = document.getElementById('progress-bar-fill');
const progressText = document.getElementById('progress-text');
const filesTableBody = document.getElementById('files-table-body');
const noFilesMessage = document.getElementById('no-files-message');
const selectAllBtn = document.getElementById('select-all-btn');
const selectAllCheckbox = document.getElementById('select-all-checkbox');
const deleteSelectedBtn = document.getElementById('delete-selected-btn');
const confirmationModal = document.getElementById('confirmation-modal');
const confirmationMessage = document.getElementById('confirmation-message');
const cancelDeleteBtn = document.getElementById('cancel-delete-btn');
const confirmDeleteBtn = document.getElementById('confirm-delete-btn');

// State
let files = [];
let selectedFiles = new Set();
let deleteCallback = null;

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    // Set up event listeners
    setupEventListeners();

    // Load files on page load
    loadFiles();
});

// Set up all event listeners
function setupEventListeners() {
    // File input change
    fileInput.addEventListener('change', handleFileInputChange);

    // Drag and drop
    dropZone.addEventListener('dragover', handleDragOver);
    dropZone.addEventListener('dragleave', handleDragLeave);
    dropZone.addEventListener('drop', handleDrop);
    dropZone.addEventListener('click', (event) => {
        // Don't trigger fileInput.click() if the click was on the label
        // as the label already has a 'for' attribute that does this
        if (event.target.tagName !== 'LABEL') {
            fileInput.click();
        }
    });

    // File selection and deletion
    selectAllBtn.addEventListener('click', toggleSelectAll);
    selectAllCheckbox.addEventListener('change', handleSelectAllCheckboxChange);
    deleteSelectedBtn.addEventListener('click', handleDeleteSelected);

    // Modal buttons
    cancelDeleteBtn.addEventListener('click', hideConfirmationModal);
    confirmDeleteBtn.addEventListener('click', handleConfirmDelete);
}

// Load files from the server
async function loadFiles() {
    try {
        const response = await fetch('/files');
        if (!response.ok) {
            throw new Error('Failed to load files');
        }

        files = await response.json();
        renderFiles();
    } catch (error) {
        console.error('Error loading files:', error);
        showError('Failed to load files. Please try again.');
    }
}

// Render files in the table
function renderFiles() {
    // Clear the table
    filesTableBody.innerHTML = '';

    // Show/hide no files message
    if (files.length === 0) {
        noFilesMessage.classList.remove('hidden');
        return;
    }

    noFilesMessage.classList.add('hidden');

    // Add files to the table
    files.forEach(file => {
        const row = document.createElement('tr');
        row.className = 'hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200';
        row.dataset.id = file.id;

        // Format file size
        const formattedSize = formatFileSize(file.size);

        // Format date
        const formattedDate = formatDate(new Date(file.uploadedAt));

        row.innerHTML = `
            <td class="p-3 border-b border-border-color dark:border-dark-border-color checkbox-column transition-colors duration-200">
                <input type="checkbox" class="file-checkbox" data-id="${file.id}">
            </td>
            <td class="p-3 border-b border-border-color dark:border-dark-border-color transition-colors duration-200">
                <span class="font-medium text-primary dark:text-dark-primary cursor-pointer hover:underline file-name transition-colors duration-200" data-id="${file.id}">${file.name}</span>
            </td>
            <td class="p-3 border-b border-border-color dark:border-dark-border-color text-text-lighter dark:text-dark-text-lighter text-sm file-size transition-colors duration-200">${formattedSize}</td>
            <td class="p-3 border-b border-border-color dark:border-dark-border-color text-text-lighter dark:text-dark-text-lighter text-sm file-date transition-colors duration-200">${formattedDate}</td>
            <td class="p-3 border-b border-border-color dark:border-dark-border-color flex gap-2 file-actions transition-colors duration-200">
                <button class="px-3 py-1 rounded bg-primary dark:bg-dark-primary text-white hover:bg-primary-hover dark:hover:bg-dark-primary-hover transition-all duration-300 action-btn download" data-id="${file.id}" title="Download">
                    Download
                </button>
                <button class="px-3 py-1 rounded bg-gray-500 dark:bg-gray-600 text-white hover:bg-gray-600 dark:hover:bg-gray-700 transition-all duration-300 action-btn view" data-id="${file.id}" title="View">
                    View
                </button>
                <button class="px-3 py-1 rounded bg-danger dark:bg-dark-danger text-white hover:bg-danger-hover dark:hover:bg-dark-danger-hover transition-all duration-300 action-btn delete" data-id="${file.id}" title="Delete">
                    Delete
                </button>
            </td>
        `;

        filesTableBody.appendChild(row);
    });

    // Add event listeners to the new elements
    addFileRowEventListeners();
}

// Add event listeners to file rows
function addFileRowEventListeners() {
    // File checkboxes
    document.querySelectorAll('.file-checkbox').forEach(checkbox => {
        checkbox.addEventListener('change', handleFileCheckboxChange);
    });

    // File names (for download)
    document.querySelectorAll('.file-name').forEach(name => {
        name.addEventListener('click', handleFileNameClick);
    });

    // Download buttons
    document.querySelectorAll('.action-btn.download').forEach(button => {
        button.addEventListener('click', handleDownloadClick);
    });

    // View buttons
    document.querySelectorAll('.action-btn.view').forEach(button => {
        button.addEventListener('click', handleViewClick);
    });

    // Delete buttons
    document.querySelectorAll('.action-btn.delete').forEach(button => {
        button.addEventListener('click', handleDeleteClick);
    });
}

// Handle file input change
function handleFileInputChange(event) {
    const files = event.target.files;
    if (files.length > 0) {
        uploadFiles(files);
    }
}

// Handle drag over
function handleDragOver(event) {
    event.preventDefault();
    event.stopPropagation();
    dropZone.classList.add('drag-over');
}

// Handle drag leave
function handleDragLeave(event) {
    event.preventDefault();
    event.stopPropagation();
    dropZone.classList.remove('drag-over');
}

// Handle drop
function handleDrop(event) {
    event.preventDefault();
    event.stopPropagation();
    dropZone.classList.remove('drag-over');

    const files = event.dataTransfer.files;
    if (files.length > 0) {
        uploadFiles(files);
    }
}

// Upload files to the server
async function uploadFiles(fileList) {
    // Show progress bar
    uploadProgress.classList.remove('hidden');
    progressBarFill.style.width = '0%';
    progressText.textContent = 'Uploading...';

    const formData = new FormData();
    for (let i = 0; i < fileList.length; i++) {
        formData.append('files', fileList[i]);
    }

    try {
        const xhr = new XMLHttpRequest();

        // Set up progress tracking
        xhr.upload.addEventListener('progress', event => {
            if (event.lengthComputable) {
                const percentComplete = Math.round((event.loaded / event.total) * 100);
                progressBarFill.style.width = `${percentComplete}%`;
                progressText.textContent = `Uploading... ${percentComplete}%`;
            }
        });

        // Set up completion handler
        xhr.addEventListener('load', () => {
            if (xhr.status >= 200 && xhr.status < 300) {
                progressBarFill.style.width = '100%';
                progressText.textContent = 'Upload complete!';

                // Parse the response
                const uploadedFiles = JSON.parse(xhr.responseText);

                // Add the new files to our list and re-render
                files = [...files, ...uploadedFiles];
                renderFiles();

                // Hide progress bar after a delay
                setTimeout(() => {
                    uploadProgress.classList.add('hidden');
                    // Reset the file input
                    fileInput.value = '';
                }, 2000);
            } else {
                handleUploadError(xhr.statusText);
            }
        });

        // Set up error handler
        xhr.addEventListener('error', () => {
            handleUploadError('Network error');
        });

        // Send the request
        xhr.open('POST', '/upload');
        xhr.send(formData);
    } catch (error) {
        handleUploadError(error.message);
    }
}

// Handle upload error
function handleUploadError(message) {
    progressText.textContent = `Upload failed: ${message}`;
    progressBarFill.style.width = '0%';

    // Hide progress bar after a delay
    setTimeout(() => {
        uploadProgress.classList.add('hidden');
        // Reset the file input
        fileInput.value = '';
    }, 3000);
}

// Handle file checkbox change
function handleFileCheckboxChange(event) {
    const fileId = event.target.dataset.id;

    if (event.target.checked) {
        selectedFiles.add(fileId);
    } else {
        selectedFiles.delete(fileId);
    }

    updateDeleteSelectedButton();
    updateSelectAllCheckbox();
}

// Handle select all checkbox change
function handleSelectAllCheckboxChange(event) {
    const checkboxes = document.querySelectorAll('.file-checkbox');

    checkboxes.forEach(checkbox => {
        checkbox.checked = event.target.checked;

        if (event.target.checked) {
            selectedFiles.add(checkbox.dataset.id);
        } else {
            selectedFiles.delete(checkbox.dataset.id);
        }
    });

    updateDeleteSelectedButton();
}

// Toggle select all
function toggleSelectAll() {
    const allSelected = selectedFiles.size === files.length;

    if (allSelected) {
        // Deselect all
        selectedFiles.clear();
    } else {
        // Select all
        files.forEach(file => {
            selectedFiles.add(file.id);
        });
    }

    // Update checkboxes
    document.querySelectorAll('.file-checkbox').forEach(checkbox => {
        checkbox.checked = !allSelected;
    });

    // Update select all checkbox
    selectAllCheckbox.checked = !allSelected;

    updateDeleteSelectedButton();
}

// Update delete selected button state
function updateDeleteSelectedButton() {
    deleteSelectedBtn.disabled = selectedFiles.size === 0;
}

// Update select all checkbox state
function updateSelectAllCheckbox() {
    const checkboxes = document.querySelectorAll('.file-checkbox');
    const checkedCount = document.querySelectorAll('.file-checkbox:checked').length;

    selectAllCheckbox.checked = checkedCount === checkboxes.length && checkboxes.length > 0;
    selectAllCheckbox.indeterminate = checkedCount > 0 && checkedCount < checkboxes.length;
}

// Handle file name click (download)
function handleFileNameClick(event) {
    const fileId = event.target.dataset.id;
    downloadFile(fileId);
}

// Handle download button click
function handleDownloadClick(event) {
    const fileId = event.target.dataset.id;
    downloadFile(fileId);
}

// Download a file
function downloadFile(fileId) {
    window.location.href = `/download/${fileId}`;
}

// Handle view button click
function handleViewClick(event) {
    const fileId = event.target.dataset.id;
    viewFile(fileId);
}

// View a file in the browser
function viewFile(fileId) {
    window.open(`/view/${fileId}`, '_blank');
}

// Handle delete button click
function handleDeleteClick(event) {
    const fileId = event.target.dataset.id;
    const file = files.find(f => f.id === fileId);

    if (file) {
        showConfirmationModal(`Are you sure you want to delete "${file.name}"?`, () => {
            deleteFile(fileId);
        });
    }
}

// Handle delete selected button click
function handleDeleteSelected() {
    if (selectedFiles.size === 0) return;

    const message = selectedFiles.size === 1
        ? 'Are you sure you want to delete the selected file?'
        : `Are you sure you want to delete ${selectedFiles.size} selected files?`;

    showConfirmationModal(message, () => {
        deleteMultipleFiles(Array.from(selectedFiles));
    });
}

// Show confirmation modal
function showConfirmationModal(message, callback) {
    confirmationMessage.textContent = message;
    deleteCallback = callback;
    confirmationModal.classList.remove('hidden');
}

// Hide confirmation modal
function hideConfirmationModal() {
    confirmationModal.classList.add('hidden');
    deleteCallback = null;
}

// Handle confirm delete
function handleConfirmDelete() {
    if (deleteCallback) {
        deleteCallback();
    }
    hideConfirmationModal();
}

// Delete a single file
async function deleteFile(fileId) {
    try {
        const formData = new FormData();
        formData.append('id', fileId);

        const response = await fetch('/delete', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Failed to delete file');
        }

        // Remove the file from our list
        files = files.filter(file => file.id !== fileId);
        selectedFiles.delete(fileId);

        // Re-render the files
        renderFiles();
        updateDeleteSelectedButton();
    } catch (error) {
        console.error('Error deleting file:', error);
        showError('Failed to delete file. Please try again.');
    }
}

// Delete multiple files
async function deleteMultipleFiles(fileIds) {
    try {
        const formData = new FormData();
        fileIds.forEach(id => {
            formData.append('ids', id);
        });

        const response = await fetch('/delete-multiple', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Failed to delete files');
        }

        // Remove the files from our list
        files = files.filter(file => !fileIds.includes(file.id));
        fileIds.forEach(id => {
            selectedFiles.delete(id);
        });

        // Re-render the files
        renderFiles();
        updateDeleteSelectedButton();
    } catch (error) {
        console.error('Error deleting files:', error);
        showError('Failed to delete files. Please try again.');
    }
}

// Show error message
function showError(message) {
    // In a real application, you might want to show a toast or alert
    alert(message);
}

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Format date
function formatDate(date) {
    const options = { 
        year: 'numeric', 
        month: 'short', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    };

    return date.toLocaleDateString(undefined, options);
}
