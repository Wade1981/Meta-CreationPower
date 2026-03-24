#include "TrayIcon.h"
#include <QIcon>
#include <QAction>

TrayIcon::TrayIcon(QObject *parent) : QObject(parent)
{
    setupTrayIcon();
    setupMenu();
}

TrayIcon::~TrayIcon()
{
    delete trayMenu;
    delete trayIcon;
}

void TrayIcon::setupTrayIcon()
{
    trayIcon = new QSystemTrayIcon(this);
    trayIcon->setIcon(QIcon("://icons/elr_icon.png"));
    trayIcon->setToolTip("ELR Desktop Assistant");
    trayIcon->show();
    
    connect(trayIcon, &QSystemTrayIcon::activated, this, &TrayIcon::onTrayIconActivated);
}

void TrayIcon::setupMenu()
{
    trayMenu = new QMenu();
    
    QAction *showWidgetAction = new QAction("显示桌面助手", this);
    QAction *hideWidgetAction = new QAction("隐藏桌面助手", this);
    QAction *exitAction = new QAction("退出", this);
    
    connect(showWidgetAction, &QAction::triggered, this, &TrayIcon::onShowWidget);
    connect(hideWidgetAction, &QAction::triggered, this, &TrayIcon::onHideWidget);
    connect(exitAction, &QAction::triggered, this, &TrayIcon::onExit);
    
    trayMenu->addAction(showWidgetAction);
    trayMenu->addAction(hideWidgetAction);
    trayMenu->addSeparator();
    trayMenu->addAction(exitAction);
    
    trayIcon->setContextMenu(trayMenu);
}

void TrayIcon::updateStatus(const QString &status)
{
    QString toolTip = "ELR Desktop Assistant\n状态: " + status;
    trayIcon->setToolTip(toolTip);
    
    // 根据状态更新图标
    if (status == "运行中") {
        trayIcon->setIcon(QIcon("://icons/elr_icon_running.png"));
    } else if (status == "停止") {
        trayIcon->setIcon(QIcon("://icons/elr_icon_stopped.png"));
    } else {
        trayIcon->setIcon(QIcon("://icons/elr_icon.png"));
    }
}

void TrayIcon::onTrayIconActivated(QSystemTrayIcon::ActivationReason reason)
{
    if (reason == QSystemTrayIcon::Trigger) {
        // 双击显示/隐藏桌面助手
        emit showDesktopWidget();
    }
}

void TrayIcon::onShowWidget()
{
    emit showDesktopWidget();
}

void TrayIcon::onHideWidget()
{
    emit hideDesktopWidget();
}

void TrayIcon::onExit()
{
    emit exitApp();
}
