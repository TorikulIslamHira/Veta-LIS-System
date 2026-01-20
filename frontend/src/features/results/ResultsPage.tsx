import React, { useState } from 'react';
import Card from '../../components/Card';
import Button from '../../components/Button';
import StatusBadge from '../../components/StatusBadge';

interface LabResult {
  id: string;
  petName: string;
  testType: string;
  testDate: string;
  status: 'normal' | 'abnormal' | 'critical';
  results: {
    parameter: string;
    value: string;
    unit: string;
    referenceRange: string;
    status: 'normal' | 'abnormal' | 'critical';
  }[];
}

const ResultsPage: React.FC = () => {
  const [results] = useState<LabResult[]>([
    {
      id: '1',
      petName: 'Max',
      testType: 'Complete Blood Count (CBC)',
      testDate: '2024-01-15 10:30 AM',
      status: 'normal',
      results: [
        { parameter: 'WBC', value: '7.5', unit: 'K/µL', referenceRange: '6.0-17.0', status: 'normal' },
        { parameter: 'RBC', value: '6.8', unit: 'M/µL', referenceRange: '5.5-8.5', status: 'normal' },
        { parameter: 'Hemoglobin', value: '15.2', unit: 'g/dL', referenceRange: '12.0-18.0', status: 'normal' },
        { parameter: 'Platelets', value: '245', unit: 'K/µL', referenceRange: '200-500', status: 'normal' },
      ],
    },
    {
      id: '2',
      petName: 'Luna',
      testType: 'Blood Chemistry',
      testDate: '2024-01-14 2:15 PM',
      status: 'abnormal',
      results: [
        { parameter: 'Glucose', value: '125', unit: 'mg/dL', referenceRange: '70-110', status: 'abnormal' },
        { parameter: 'BUN', value: '22', unit: 'mg/dL', referenceRange: '7-27', status: 'normal' },
        { parameter: 'Creatinine', value: '1.2', unit: 'mg/dL', referenceRange: '0.5-1.5', status: 'normal' },
        { parameter: 'ALT', value: '45', unit: 'U/L', referenceRange: '10-125', status: 'normal' },
      ],
    },
    {
      id: '3',
      petName: 'Charlie',
      testType: 'Urinalysis',
      testDate: '2024-01-12 9:00 AM',
      status: 'critical',
      results: [
        { parameter: 'pH', value: '8.5', unit: '', referenceRange: '5.5-7.0', status: 'critical' },
        { parameter: 'Protein', value: 'Positive', unit: '', referenceRange: 'Negative', status: 'abnormal' },
        { parameter: 'Glucose', value: 'Negative', unit: '', referenceRange: 'Negative', status: 'normal' },
      ],
    },
  ]);

  const [selectedResult, setSelectedResult] = useState<LabResult | null>(null);
  const [filter, setFilter] = useState<'all' | 'normal' | 'abnormal' | 'critical'>('all');

  const filteredResults = filter === 'all' 
    ? results 
    : results.filter(r => r.status === filter);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Test Results</h2>
          <p className="text-gray-600 mt-1">View and analyze laboratory test results</p>
        </div>
        <Button>📥 Export Results</Button>
      </div>

      {/* Filter Tabs */}
      <Card>
        <div className="flex gap-2">
          {(['all', 'normal', 'abnormal', 'critical'] as const).map((status) => (
            <button
              key={status}
              onClick={() => setFilter(status)}
              className={`px-4 py-2 rounded-lg font-medium transition-all ${
                filter === status
                  ? 'bg-primary text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {status.charAt(0).toUpperCase() + status.slice(1)}
              {status !== 'all' && (
                <span className="ml-2 text-xs opacity-75">
                  ({results.filter(r => r.status === status).length})
                </span>
              )}
            </button>
          ))}
        </div>
      </Card>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Results List */}
        <div className="lg:col-span-2">
          <Card>
            <div className="space-y-3">
              {filteredResults.length === 0 ? (
                <div className="text-center py-12">
                  <div className="text-6xl mb-4">📋</div>
                  <p className="text-gray-500">No results found</p>
                </div>
              ) : (
                filteredResults.map((result) => (
                  <div
                    key={result.id}
                    onClick={() => setSelectedResult(result)}
                    className={`p-4 border rounded-lg cursor-pointer transition-all ${
                      selectedResult?.id === result.id
                        ? 'border-primary bg-blue-50'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <div className="flex items-start justify-between mb-2">
                      <div>
                        <h4 className="font-semibold text-gray-900">{result.testType}</h4>
                        <p className="text-sm text-gray-600">Pet: {result.petName}</p>
                        <p className="text-xs text-gray-500 mt-1">{result.testDate}</p>
                      </div>
                      <StatusBadge status={result.status} />
                    </div>
                    <div className="text-sm text-gray-600">
                      {result.results.length} parameters tested
                    </div>
                  </div>
                ))
              )}
            </div>
          </Card>
        </div>

        {/* Result Details */}
        <div>
          <Card title="Result Details">
            {selectedResult ? (
              <div className="space-y-4">
                <div>
                  <h3 className="font-bold text-lg">{selectedResult.testType}</h3>
                  <p className="text-sm text-gray-600 mt-1">Pet: {selectedResult.petName}</p>
                  <p className="text-xs text-gray-500">{selectedResult.testDate}</p>
                  <div className="mt-2">
                    <StatusBadge status={selectedResult.status} />
                  </div>
                </div>

                <div className="border-t pt-4">
                  <h4 className="font-semibold mb-3">Test Parameters</h4>
                  <div className="space-y-3">
                    {selectedResult.results.map((param, index) => (
                      <div key={index} className="bg-gray-50 p-3 rounded-lg">
                        <div className="flex justify-between items-start mb-1">
                          <span className="font-medium text-sm">{param.parameter}</span>
                          <StatusBadge status={param.status} />
                        </div>
                        <div className="text-lg font-bold text-gray-900">
                          {param.value} {param.unit}
                        </div>
                        <div className="text-xs text-gray-500 mt-1">
                          Reference: {param.referenceRange}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="border-t pt-4 space-y-2">
                  <Button variant="primary" className="w-full">
                    📄 Download PDF
                  </Button>
                  <Button variant="secondary" className="w-full">
                    📧 Email to Owner
                  </Button>
                  <Button variant="secondary" className="w-full">
                    📋 Print Report
                  </Button>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                Select a result to view details
              </div>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
};

export default ResultsPage;
