import React, { useEffect, useState } from 'react';
import { GetDashboardStats, GetRecentPets } from '../../../wailsjs/go/app/App';
import { EventsOn } from '../../../wailsjs/runtime/runtime';
import { core } from '../../../wailsjs/go/models';
import Card from '../../components/Card';
import StatusBadge from '../../components/StatusBadge';
import Button from '../../components/Button';

interface Activity {
  time: string;
  action: string;
  type: 'success' | 'info' | 'warning' | 'error';
}

export const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    totalMachines: 1,
    activeMachines: 1,
    totalTests: 0,
    todayTests: 0,
  });
  const [activities, setActivities] = useState<Activity[]>([]);

  useEffect(() => {
    const initDashboard = async () => {
      try {
        // Load dashboard statistics from backend
        const dashboardStats = await GetDashboardStats();
        setStats({
          totalMachines: 1,
          activeMachines: 1,
          totalTests: dashboardStats.total_tests || 0,
          todayTests: dashboardStats.pending || 0,
        });

        // Load recent pets to populate activity feed
        const recentPets = await GetRecentPets();
        const petActivities: Activity[] = recentPets.slice(0, 4).map((pet: core.Pet, index: number) => ({
          time: `${(index + 1) * 15} min ago`,
          action: `Registered: ${pet.name} (${pet.species})`,
          type: 'info' as const,
        }));
        setActivities(petActivities);

      } catch (error) {
        console.error('Failed to initialize dashboard:', error);
      } finally {
        setLoading(false);
      }
    };

    initDashboard();

    // Set up real-time event listener for new lab results
    EventsOn('new_result', (result: core.LabResult) => {
      console.log('New result received:', result);
      
      // Determine result type based on status or value
      const resultType: Activity['type'] = 
        result.status === 'critical' ? 'warning' : 'success';

      // Add new activity to the feed
      const newActivity: Activity = {
        time: 'Just now',
        action: `New Result: ${result.test_name} = ${result.result_value} ${result.unit || ''}`,
        type: resultType,
      };

      setActivities(prev => [newActivity, ...prev.slice(0, 3)]);

      // Update stats
      setStats(prev => ({
        ...prev,
        totalTests: prev.totalTests + 1,
        todayTests: prev.todayTests + 1,
      }));
    });

    // Cleanup event listener on unmount
    return () => {
      // Wails handles cleanup automatically
    };
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-xl text-gray-600">Loading machines...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Machines</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">{stats.totalMachines}</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center text-2xl">
              🔬
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Active Machines</p>
              <p className="text-3xl font-bold text-green-600 mt-2">{stats.activeMachines}</p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center text-2xl">
              ✓
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Tests</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">{stats.totalTests}</p>
            </div>
            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center text-2xl">
              📊
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Today's Tests</p>
              <p className="text-3xl font-bold text-blue-600 mt-2">{stats.todayTests}</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center text-2xl">
              📅
            </div>
          </div>
        </Card>
      </div>

      {/* Lab Machines */}
      <Card title="Lab Machines" actions={
        <Button onClick={() => window.location.reload()} variant="secondary" size="sm">
          🔄 Refresh
        </Button>
      }>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {/* Static machine card - backend handles connection internally */}
          <div className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div className="flex items-start justify-between mb-3">
              <div>
                <h4 className="font-semibold text-gray-900">Scitek</h4>
                <p className="text-sm text-gray-500">Model: SC-2000</p>
              </div>
              <StatusBadge status="connected" />
            </div>

            <div className="space-y-2 mb-4">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Port:</span>
                <span className="font-medium">COM3</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Type:</span>
                <span className="font-medium">Serial</span>
              </div>
            </div>

            <div className="flex gap-2">
              <Button
                onClick={() => console.log('Machine configured in backend')}
                variant="success"
                size="sm"
                className="flex-1"
                disabled
              >
                Active
              </Button>
              <Button variant="secondary" size="sm">
                ⚙️
              </Button>
            </div>
          </div>
        </div>
      </Card>

      {/* Recent Activity */}
      <Card title="Recent Activity">
        <div className="space-y-3">
          {activities.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              No recent activity
            </div>
          ) : (
            activities.map((activity, index) => (
              <div key={index} className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
                <div className={`w-2 h-2 rounded-full mt-2 ${
                  activity.type === 'success' ? 'bg-green-500' :
                  activity.type === 'info' ? 'bg-blue-500' :
                  activity.type === 'warning' ? 'bg-yellow-500' : 'bg-red-500'
                }`}></div>
                <div className="flex-1">
                  <p className="text-sm text-gray-900">{activity.action}</p>
                  <p className="text-xs text-gray-500 mt-1">{activity.time}</p>
                </div>
              </div>
            ))
          )}
        </div>
      </Card>
    </div>
  );
};

export default Dashboard;

