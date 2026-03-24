#include "ELRDesktopAssistant.h"

ELRDesktopAssistant::ELRDesktopAssistant(QWidget *parent)
    : QMainWindow(parent)
{
    initializeComponents();
    setupConnections();
    elrClient->startMonitoring();
}

ELRDesktopAssistant::~ELRDesktopAssistant()
{
    delete trayIcon;
    delete desktopWidget;
    delete elrClient;
}

void ELRDesktopAssistant::initializeComponents()
{
    // 创建系统托盘图标
    trayIcon = new TrayIcon(this);
    
    // 创建桌面助手 widget
    desktopWidget = new DesktopWidget(this);
    
    // 创建 ELR 客户端
    elrClient = new ELRClient(this);
    
    // 隐藏主窗口
    setVisible(false);
}

void ELRDesktopAssistant::setupConnections()
{
    // 连接 ELR 客户端的信号
    connect(elrClient, &ELRClient::statusUpdated, this, &ELRDesktopAssistant::onELRStatusUpdated);
    connect(elrClient, &ELRClient::containersUpdated, this, &ELRDesktopAssistant::onELRContainersUpdated);
    
    // 连接托盘图标的信号
    connect(trayIcon, &TrayIcon::showDesktopWidget, desktopWidget, &DesktopWidget::show);
    connect(trayIcon, &TrayIcon::hideDesktopWidget, desktopWidget, &DesktopWidget::hide);
    connect(trayIcon, &TrayIcon::exitApp, qApp, &QApplication::quit);
}

void ELRDesktopAssistant::onELRStatusUpdated(const QString &status)
{
    trayIcon->updateStatus(status);
    desktopWidget->updateStatus(status);
}

void ELRDesktopAssistant::onELRContainersUpdated(const QList<QMap<QString, QString>> &containers)
{
    desktopWidget->updateContainers(containers);
}
