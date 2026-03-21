"""Archive service module for managing EL-CSCC Archive."""

import json
import os
from datetime import datetime


class ArchiveService:
    """Service for managing EL-CSCC Archive."""

    def __init__(self, archive_file):
        """
        Initialize ArchiveService with archive file path.
        
        Args:
            archive_file (str): Path to archive JSON file
        """
        self.archive_file = archive_file
        self.archives = self._load_archives()

    def _load_archives(self):
        """
        Load archives from JSON file.
        
        Returns:
            dict: Archives dictionary
        """
        if os.path.exists(self.archive_file):
            try:
                with open(self.archive_file, 'r', encoding='utf-8') as f:
                    return json.load(f)
            except Exception as e:
                print(f"Error loading archives: {e}")
                return {}
        return {}

    def _save_archives(self):
        """
        Save archives to JSON file.
        """
        try:
            with open(self.archive_file, 'w', encoding='utf-8') as f:
                json.dump(self.archives, f, ensure_ascii=False, indent=2)
            return True
        except Exception as e:
            print(f"Error saving archives: {e}")
            return False

    def get_archive(self, archive_id):
        """
        Get archive by ID.
        
        Args:
            archive_id (str): Archive ID
            
        Returns:
            dict: Archive information or None if not found
        """
        return self.archives.get(archive_id)

    def list_archives(self):
        """
        List all archives.
        
        Returns:
            list: List of archive information
        """
        return list(self.archives.values())

    def add_archive(self, archive_data):
        """
        Add new archive.
        
        Args:
            archive_data (dict): Archive data
            
        Returns:
            dict: Added archive with ID
        """
        # Generate archive ID if not provided
        if 'archive_id' not in archive_data:
            archive_data['archive_id'] = f"ARC-{datetime.now().strftime('%Y%m%d%H%M%S')}"
        
        # Add timestamp if not provided
        if 'creation_time' not in archive_data:
            archive_data['creation_time'] = datetime.now().isoformat()
        
        # Add to archives
        self.archives[archive_data['archive_id']] = archive_data
        
        # Save to file
        self._save_archives()
        
        return archive_data

    def update_archive(self, archive_id, archive_data):
        """
        Update existing archive.
        
        Args:
            archive_id (str): Archive ID
            archive_data (dict): Updated archive data
            
        Returns:
            dict: Updated archive or None if not found
        """
        if archive_id in self.archives:
            # Update archive
            self.archives[archive_id].update(archive_data)
            
            # Add update timestamp
            self.archives[archive_id]['last_updated'] = datetime.now().isoformat()
            
            # Save to file
            self._save_archives()
            
            return self.archives[archive_id]
        return None

    def delete_archive(self, archive_id):
        """
        Delete archive by ID.
        
        Args:
            archive_id (str): Archive ID
            
        Returns:
            bool: True if deleted, False if not found
        """
        if archive_id in self.archives:
            del self.archives[archive_id]
            self._save_archives()
            return True
        return False

    def search_archives(self, query):
        """
        Search archives by query.
        
        Args:
            query (str): Search query
            
        Returns:
            list: List of matching archives
        """
        results = []
        query_lower = query.lower()
        
        for archive in self.archives.values():
            # Search in multiple fields
            for field in ['title', 'description', 'carbon_partner', 'silicon_partner']:
                if field in archive and query_lower in str(archive[field]).lower():
                    results.append(archive)
                    break
        
        return results
