// Example usage of the MachineIDService in TypeScript
import { MachineIDService } from './frontend/bindings/github.com/kawai-network/veridium';

// Get the raw machine ID (consider this confidential)
async function getMachineID() {
  try {
    const machineID = await MachineIDService.GetID();
    console.log('Machine ID:', machineID);
    return machineID;
  } catch (error) {
    console.error('Failed to get machine ID:', error);
  }
}

// Get a protected/hashed machine ID using an application-specific key
// This is recommended for sharing machine IDs securely
async function getProtectedMachineID() {
  try {
    const appID = 'your-app-identifier'; // Use a unique identifier for your app
    const protectedID = await MachineIDService.GetProtectedID(appID);
    console.log('Protected Machine ID:', protectedID);
    return protectedID;
  } catch (error) {
    console.error('Failed to get protected machine ID:', error);
  }
}

// Example usage
export async function example() {
  await getMachineID();
  await getProtectedMachineID();
}
