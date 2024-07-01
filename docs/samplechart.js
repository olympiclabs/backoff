

const ctx = document.getElementById('myChart');

new Chart(ctx, {
  type: 'line',

  data: {
    datasets: [{
      data: [
        {x: '1', y: 16.094},
        {x: '2', y: 32.189},
        {x: '3', y: 38.067},
        {x: '4', y: 41.744},
        {x: '5', y: 44.427},
        {x: '6', y: 46.540},
        {x: '7', y: 48.283},
        {x: '8', y: 49.767},
        {x: '9', y: 51.059},
        {x: '10', y: 52.204},
        {x: '11', y: 53.230},
        {x: '12', y: 54.161},
        {x: '13', y: 55.013},
        {x: '14', y: 55.797},
        {x: '15', y: 56.525},
        {x: '16', y: 57.203},
        {x: '17', y: 57.838},
        {x: '18', y: 58.435},
        {x: '19', y: 58.999},
        {x: '20', y: 59.532},
        {x: '21', y: 60.039},
      ],
      borderWidth: 1
    }]
  },

  options: {
    plugins: { legend: { display: false} },
    cubicInterpolationMode: 'monotone',
    scales: {
      y: {
        beginAtZero: true
      }
    }
  }
});
