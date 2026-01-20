import React, { useState } from 'react';
import Card from '../../components/Card';
import Button from '../../components/Button';

interface Pet {
  id: string;
  name: string;
  species: string;
  breed: string;
  age: number;
  owner: string;
  lastVisit: string;
}

const PetsPage: React.FC = () => {
  const [pets] = useState<Pet[]>([
    {
      id: '1',
      name: 'Max',
      species: 'Dog',
      breed: 'Golden Retriever',
      age: 5,
      owner: 'John Smith',
      lastVisit: '2026-01-18',
    },
    {
      id: '2',
      name: 'Luna',
      species: 'Cat',
      breed: 'Persian',
      age: 3,
      owner: 'Sarah Johnson',
      lastVisit: '2026-01-17',
    },
    {
      id: '3',
      name: 'Charlie',
      species: 'Dog',
      breed: 'Labrador',
      age: 7,
      owner: 'Mike Wilson',
      lastVisit: '2026-01-15',
    },
  ]);

  const [searchTerm, setSearchTerm] = useState('');
  const [selectedPet, setSelectedPet] = useState<Pet | null>(null);

  const filteredPets = pets.filter(
    (pet) =>
      pet.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      pet.owner.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Pets & Patients</h2>
          <p className="text-gray-600 mt-1">Manage patient information and medical records</p>
        </div>
        <Button>+ Add New Pet</Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Pets List */}
        <div className="lg:col-span-2">
          <Card>
            <div className="mb-4">
              <input
                type="text"
                placeholder="Search by pet name or owner..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
              />
            </div>

            <div className="space-y-3">
              {filteredPets.length === 0 ? (
                <div className="text-center py-12">
                  <div className="text-6xl mb-4">🐕</div>
                  <p className="text-gray-500 mb-4">No pets found</p>
                  <Button>Add Your First Pet</Button>
                </div>
              ) : (
                filteredPets.map((pet) => (
                  <div
                    key={pet.id}
                    onClick={() => setSelectedPet(pet)}
                    className={`p-4 border rounded-lg cursor-pointer transition-all ${
                      selectedPet?.id === pet.id
                        ? 'border-primary bg-blue-50'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <div className="w-14 h-14 bg-gradient-to-br from-blue-100 to-blue-200 rounded-full flex items-center justify-center text-3xl">
                          {pet.species === 'Dog' ? '🐕' : '🐈'}
                        </div>
                        <div>
                          <h4 className="font-semibold text-gray-900">{pet.name}</h4>
                          <p className="text-sm text-gray-600">
                            {pet.breed} • {pet.age} years old
                          </p>
                          <p className="text-xs text-gray-500 mt-1">
                            Owner: {pet.owner}
                          </p>
                        </div>
                      </div>

                      <div className="text-right">
                        <p className="text-sm text-gray-600">Last Visit</p>
                        <p className="text-sm font-medium">{pet.lastVisit}</p>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </Card>
        </div>

        {/* Pet Details */}
        <div>
          <Card title="Pet Details">
            {selectedPet ? (
              <div className="space-y-4">
                <div className="text-center">
                  <div className="w-20 h-20 bg-gradient-to-br from-blue-100 to-blue-200 rounded-full flex items-center justify-center text-5xl mx-auto mb-3">
                    {selectedPet.species === 'Dog' ? '🐕' : '🐈'}
                  </div>
                  <h3 className="text-xl font-bold">{selectedPet.name}</h3>
                  <p className="text-gray-600">{selectedPet.breed}</p>
                </div>

                <div className="space-y-3 pt-4">
                  <div>
                    <label className="text-sm font-medium text-gray-600">Species</label>
                    <p className="text-gray-900 mt-1">{selectedPet.species}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-600">Age</label>
                    <p className="text-gray-900 mt-1">{selectedPet.age} years</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-600">Owner</label>
                    <p className="text-gray-900 mt-1">{selectedPet.owner}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-600">Last Visit</label>
                    <p className="text-gray-900 mt-1">{selectedPet.lastVisit}</p>
                  </div>
                </div>

                <div className="pt-4 border-t space-y-2">
                  <Button variant="primary" className="w-full">
                    📋 View Medical Records
                  </Button>
                  <Button variant="secondary" className="w-full">
                    ✏️ Edit Information
                  </Button>
                  <Button variant="secondary" className="w-full">
                    🔬 New Test
                  </Button>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                Select a pet to view details
              </div>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
};

export default PetsPage;
