package store

func (f FolderFilter) Matches(chat Chat) bool {
	// Excluded peers never match
	for _, id := range f.ExcludePeers {
		if id == chat.ID {
			return false
		}
	}

	// Explicitly included/pinned peers always match (bypass category and exclusion flags)
	for _, id := range f.IncludePeers {
		if id == chat.ID {
			return true
		}
	}
	for _, id := range f.PinnedPeers {
		if id == chat.ID {
			return true
		}
	}

	// Must match at least one category flag
	categoryMatched := false
	if f.Contacts && chat.IsContact && !chat.IsBot {
		categoryMatched = true
	}
	if f.NonContacts && chat.Peer.IsUser() && !chat.IsContact && !chat.IsBot {
		categoryMatched = true
	}
	if f.Groups && chat.Peer.IsGroup() {
		categoryMatched = true
	}
	if f.Broadcasts && chat.Peer.IsChannel() {
		categoryMatched = true
	}
	if f.Bots && chat.IsBot {
		categoryMatched = true
	}
	if !categoryMatched {
		return false
	}

	// Apply exclusion flags (category matches only; explicit peers bypass these above)
	if f.ExcludeRead && chat.UnreadCount == 0 {
		return false
	}
	if f.ExcludeMuted && chat.IsMuted {
		return false
	}
	return true
}
