
let liveUpdate = false

const products = [['BTC-USD', '#2596be'], ['ETH-USD', '#6c25be'], ['ETH-BTC', '#bea925']]

const chartNodes = products.map(([product]) => document.getElementById(`chart-${product}`))
const charts = {}

products.forEach(([product, color], idx) => {
    const datasets = [{
        label: product,
        borderColor: color,
        backgroundColor: color,
        borderWidth: 1,
        tension: 0.1,
        data: []
    }]

    const ctx = chartNodes[idx].getContext('2d')
    charts[product] = new Chart(ctx, {
        type: 'line',
        data: {
            datasets
        },
        options: {
            plugins: {
                legend: {
                    position: 'top',
                }
            },
            tooltips: {
                mode: 'index',
                intersect: false
            },
            hover: {
                mode: 'index',
                intersect: false
            },
            elements: {
                point: {
                    radius: 0
                }
            },
            scales: {
                x: {
                    beginAtZero: false,
                    ticks: {
                        display: false
                    }
                },
                y: {
                    beginAtZero: false
                },
            }
        },
    })
})

const socket = new WebSocket("ws://localhost:8080/subscribe")

socket.onmessage = function (event) {
    const current = JSON.parse(event.data)
    // console.log(current)

    Object.keys(current).forEach(product => {
        charts[product].data.datasets[0].data.push({
            x: new Date().toLocaleString(),
            y: current[product]
        })
        if (liveUpdate) {
            charts[product].update();
        }
    })
}

const toggleBtn = document.getElementById('toggle')
toggleBtn.addEventListener('click', () => {
    liveUpdate = !liveUpdate
    toggleBtn.innerText = liveUpdate ? 'Stop' : 'Start'
    products.forEach(([product]) => {
        charts[product].update();
    })
})
toggleBtn.click()