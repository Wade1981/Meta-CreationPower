// AIAgent APP - 用户交互应用

class AgentApp {
  constructor(config) {
    this.config = config || {};
    this.clients = new Map();
    this.conversations = new Map();
    this.notifications = new Map();
  }

  async initialize() {
    console.log('Initializing AIAgent APP...');
    // 初始化APP组件
    this.chatService = new ChatService();
    this.notificationService = new NotificationService();
    this.userManager = new UserManager();
    this.authService = new AuthService();
    this.uiManager = new UIManager();

    console.log('AIAgent APP initialized successfully');
    return true;
  }

  async registerClient(clientId, clientInfo) {
    console.log(`Registering client: ${clientId}`);
    const client = {
      id: clientId,
      info: clientInfo,
      lastActive: new Date(),
      status: 'online'
    };
    this.clients.set(clientId, client);
    return { success: true, clientId };
  }

  async sendMessage(senderId, recipientId, content) {
    console.log(`Sending message from ${senderId} to ${recipientId}`);
    const message = await this.chatService.sendMessage({
      senderId,
      recipientId,
      content,
      timestamp: new Date()
    });

    // 确保会话存在
    const conversationId = this.getConversationId(senderId, recipientId);
    await this.ensureConversation(conversationId, senderId, recipientId);

    // 添加消息到会话
    const conversation = this.conversations.get(conversationId);
    conversation.messages.push(message);

    return message;
  }

  async createGroup(groupId, groupName, members) {
    console.log(`Creating group: ${groupId}`);
    const group = {
      id: groupId,
      name: groupName,
      members,
      created_at: new Date(),
      messages: []
    };
    this.conversations.set(groupId, group);
    return { success: true, groupId };
  }

  async sendGroupMessage(senderId, groupId, content) {
    console.log(`Sending group message from ${senderId} to ${groupId}`);
    const group = this.conversations.get(groupId);
    if (!group) {
      throw new Error(`Group ${groupId} not found`);
    }

    const message = await this.chatService.sendGroupMessage({
      senderId,
      groupId,
      content,
      timestamp: new Date()
    });

    group.messages.push(message);
    return message;
  }

  async sendNotification(recipientId, notificationInfo) {
    console.log(`Sending notification to ${recipientId}`);
    const notification = await this.notificationService.createNotification({
      recipientId,
      ...notificationInfo,
      timestamp: new Date()
    });

    if (!this.notifications.has(recipientId)) {
      this.notifications.set(recipientId, []);
    }
    this.notifications.get(recipientId).push(notification);

    return notification;
  }

  async getConversations(clientId) {
    console.log(`Getting conversations for client: ${clientId}`);
    const userConversations = [];

    for (const [conversationId, conversation] of this.conversations) {
      // 检查客户端是否是会话的参与者
      if (this.isParticipant(conversation, clientId)) {
        userConversations.push({
          id: conversationId,
          name: conversation.name || this.getConversationName(conversation, clientId),
          lastMessage: conversation.messages[conversation.messages.length - 1],
          messageCount: conversation.messages.length
        });
      }
    }

    return userConversations;
  }

  async getMessages(conversationId, limit = 50, offset = 0) {
    console.log(`Getting messages for conversation: ${conversationId}`);
    const conversation = this.conversations.get(conversationId);
    if (!conversation) {
      throw new Error(`Conversation ${conversationId} not found`);
    }

    return conversation.messages.slice(offset, offset + limit);
  }

  async getNotifications(clientId) {
    console.log(`Getting notifications for client: ${clientId}`);
    return this.notifications.get(clientId) || [];
  }

  async markAsRead(notificationId, clientId) {
    console.log(`Marking notification ${notificationId} as read`);
    const notifications = this.notifications.get(clientId);
    if (!notifications) {
      return { success: false, message: 'No notifications found' };
    }

    const notification = notifications.find(n => n.id === notificationId);
    if (notification) {
      notification.read = true;
      notification.read_at = new Date();
    }

    return { success: true };
  }

  // 辅助方法
  getConversationId(userId1, userId2) {
    const ids = [userId1, userId2].sort();
    return `conv_${ids.join('_')}`;
  }

  async ensureConversation(conversationId, userId1, userId2) {
    if (!this.conversations.has(conversationId)) {
      const conversation = {
        id: conversationId,
        participants: [userId1, userId2],
        messages: [],
        created_at: new Date()
      };
      this.conversations.set(conversationId, conversation);
    }
  }

  isParticipant(conversation, clientId) {
    if (conversation.participants) {
      return conversation.participants.includes(clientId);
    }
    return false;
  }

  getConversationName(conversation, clientId) {
    if (conversation.participants) {
      const otherParticipant = conversation.participants.find(id => id !== clientId);
      return otherParticipant || 'Unknown';
    }
    return 'Unknown Conversation';
  }

  async shutdown() {
    console.log('Shutting down AIAgent APP...');
    // 关闭所有服务
    console.log('AIAgent APP shutdown successfully');
    return true;
  }
}

class ChatService {
  async sendMessage(message) {
    // 发送消息
    return {
      id: `msg_${Date.now()}`,
      ...message,
      status: 'sent'
    };
  }

  async sendGroupMessage(message) {
    // 发送群消息
    return {
      id: `msg_${Date.now()}`,
      ...message,
      status: 'sent'
    };
  }
}

class NotificationService {
  async createNotification(notification) {
    // 创建通知
    return {
      id: `notif_${Date.now()}`,
      ...notification,
      read: false
    };
  }
}

class UserManager {
  async getUser(userId) {
    // 获取用户信息
    return {
      id: userId,
      name: `User ${userId}`,
      avatar: `https://avatar.example.com/${userId}`
    };
  }
}

class AuthService {
  async authenticate(credentials) {
    // 认证用户
    return {
      token: `token_${Date.now()}`,
      expires_at: new Date(Date.now() + 86400000) // 24小时
    };
  }
}

class UIManager {
  getUIConfig(clientId) {
    // 获取UI配置
    return {
      theme: 'light',
      language: 'zh-CN',
      preferences: {}
    };
  }
}

module.exports = {
  AgentApp,
  ChatService,
  NotificationService,
  UserManager,
  AuthService,
  UIManager
};