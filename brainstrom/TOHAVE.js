var data = {
    policyInfo:{
        policyNumber: 'PN123456',
        policyHolder: 'John Doe',
        coverageType: 'Comprehensive',
        startDate: '2023-01-01',
        endDate: '2024-01-01'
    },
    customerInfo:{
        name: 'John Doe',
        age:  35,
        address: '123 Main St, Anytown, USA',
        contactNumber: '555-1234',
        email: '',
    },
    claimHistory:[
        {
            claimId: 'C1001',
            date: '2023-06-15',
            amount: 5000,
            status: 'Approved',
            report:{
                reportId: 'R2001',
                description: 'Rear-end collision at traffic light',
                filedBy: 'John Doe',
                filedDate: '2023-06-16'
            }
        },
        {
            claimId: 'C1002',
            date: '2023-09-10',
            amount: 1500,
            status: 'Pending',
            report:{
                reportId: 'R2002',
                description: 'Hail damage to roof',
                filedBy: 'John Doe',
                filedDate: '2023-09-11'
            }
        }
    ],
    coverage:[
        { type: 'Collision', limit: 10000 },
        { type: 'Liability', limit: 50000 },
        { type: 'Comprehensive', limit: 15000 }
    ]
}
// Process coverage limit based on age and claim history
if(data.customerInfo.age < 25) {
    data.coverage.forEach(cov => {
        cov.limit *= 0.8; // Reduce limit by 20% for young drivers
    });
}

let totalClaims = data.claimHistory.reduce((sum, claim) => sum + claim.amount, 0);
if(totalClaims > 10000) {
    data.coverage.forEach(cov => {
        cov.limit *= 0.9; // Reduce limit by 10% if total claims exceed 10,000
    });
}