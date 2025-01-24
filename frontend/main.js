document.addEventListener('DOMContentLoaded', () => {
    fetchScans();

    document.querySelector('.sidebar').addEventListener('click', (event) => {
        if (event.target.classList.contains('scan-line')) {
            const scanId = event.target.dataset.scanId;
            fetchScanDetails(scanId);
        }
    });

    document.querySelector('.content').addEventListener('click', (event) => {
        if (event.target.classList.contains('domain-item')) {
            const credentialsDiv = event.target.querySelector('.credentials');
            credentialsDiv.style.display = credentialsDiv.style.display === 'none' ? 'block' : 'none';
        }
    });

    document.getElementById('domainForm').addEventListener('submit', (event) => {
        event.preventDefault();
        const textarea = document.getElementById('domainTextarea');
        const domains = textarea.value.split('\n').map(domain => domain.trim()).filter(domain => domain);
        sendDomains(domains);
    });
});

let cspChartInstance = null;
let xFrameChartInstance = null;

function categorizeScanDate(date) {
    const today = new Date();
    const scanDate = new Date(date);
    const diffInDays = Math.floor((today - scanDate) / (1000 * 60 * 60 * 24));

    if (diffInDays === 0) {
        return 'Today';
    } else if (diffInDays === 1) {
        return 'Yesterday';
    } else if (diffInDays <= 7) {
        return 'Previous 7 days';
    } else if (diffInDays <= 30) {
        return 'Previous month';
    } else {
        return 'Earlier';
    }
}

function formatDuration(duration) {
    const seconds = parseFloat(duration);
    if (seconds < 60) {
        return `${seconds.toFixed(2)} seconds`;
    } else if (seconds < 3600) {
        const minutes = seconds / 60;
        return `${minutes.toFixed(2)} minutes`;
    } else {
        const hours = seconds / 3600;
        return `${hours.toFixed(2)} hours`;
    }
}

function fetchScans() {
    fetch('http://localhost:8080/scans')
        .then(response => response.json())
        .then(data => {
            const sidebar = document.querySelector('.sidebar');
            sidebar.innerHTML = '';

            const categories = {};
            data.forEach(scan => {
                const category = categorizeScanDate(scan.start_time);
                if (!categories[category]) {
                    categories[category] = [];
                }
                categories[category].push(scan);
            });

            for (const category in categories) {
                const dateHeader = document.createElement('h3');
                dateHeader.textContent = category;
                sidebar.appendChild(dateHeader);

                categories[category].forEach(scan => {
                    const scanLine = document.createElement('div');
                    scanLine.classList.add('scan-line');
                    scanLine.dataset.scanId = scan.id;
                    scanLine.textContent = `${scan.domain_count} domains, ${formatDuration(scan.duration)}`;
                    sidebar.appendChild(scanLine);
                });
            }
        });
}

function fetchScanDetails(scanId) {
    fetch(`http://localhost:8080/scan/${scanId}`)
        .then(response => response.json())
        .then(data => {
            const scanInfo = document.querySelector('.scan-info');
            const domainList = document.querySelector('.domain-list');
            scanInfo.innerHTML = '';
            domainList.innerHTML = '';

            let cspTrue = 0;
            let xFrameTrue = 0;

            data.forEach(domain => {
                if (domain.csp) cspTrue++;
                if (domain.xframe) xFrameTrue++;

                const domainItem = document.createElement('div');
                domainItem.classList.add('domain-item');
                domainItem.innerHTML = `
                    <p>Domain: ${domain.domain}</p>
                    <p>CSP: ${domain.csp}</p>
                    <p>X-Frame: ${domain.xframe}</p>
                    <p>Credentials: ${Object.keys(domain.credentials).length}</p>
                    <div class="credentials">
                        ${Object.entries(domain.credentials).map(([key, value]) => `<p>${key}: ${value}</p>`).join('')}
                    </div>
                `;
                domainList.appendChild(domainItem);
            });

            const cspChartCtx = document.getElementById('cspChart').getContext('2d');
            const xFrameChartCtx = document.getElementById('xFrameChart').getContext('2d');

            // Destroy existing chart instances if they exist
            if (cspChartInstance) {
                cspChartInstance.destroy();
            }
            if (xFrameChartInstance) {
                xFrameChartInstance.destroy();
            }

            // Create new chart instances
            cspChartInstance = new Chart(cspChartCtx, {
                type: 'pie',
                data: {
                    labels: ['CSP True', 'CSP False'],
                    datasets: [{
                        data: [cspTrue, data.length - cspTrue],
                        backgroundColor: ['#4CAF50', '#FF0000']
                    }]
                }
            });

            xFrameChartInstance = new Chart(xFrameChartCtx, {
                type: 'pie',
                data: {
                    labels: ['X-Frame True', 'X-Frame False'],
                    datasets: [{
                        data: [xFrameTrue, data.length - xFrameTrue],
                        backgroundColor: ['#4CAF50', '#FF0000']
                    }]
                }
            });
        });
}

function sendDomains(domains) {
    fetch('http://localhost:8080/scan', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ domains })
    })
        .then(response => response.json())
        .then(data => {
            console.log(data);
            fetchScans(); // Update the scans sidebar
        })
        .catch(error => console.error('Error:', error));
}
