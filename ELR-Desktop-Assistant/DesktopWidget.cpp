#include "DesktopWidget.h"
#include <QHBoxLayout>
#include <QMouseEvent>
#include <QListWidgetItem>

DesktopWidget::DesktopWidget(QWidget *parent) : QWidget(parent)
{
    setupUI();
}

DesktopWidget::~DesktopWidget()
{
}

void DesktopWidget::setupUI()
{
    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowStaysOnTopHint | Qt::Tool);
    setAttribute(Qt::WA_TranslucentBackground);
    setFixedSize(300, 400);
    
    QWidget *contentWidget = new QWidget(this);
    contentWidget->setStyleSheet(
        "QWidget { background-color: rgba(255, 255, 255, 0.9); border-radius: 10px; }"
    );
    
    QVBoxLayout *mainLayout = new QVBoxLayout(contentWidget);
    
    // 标题
    QLabel *titleLabel = new QLabel("ELR 桌面助手");
    titleLabel->setAlignment(Qt::AlignCenter);
    titleLabel->setStyleSheet("font-size: 16px; font-weight: bold;");
    mainLayout->addWidget(titleLabel);
    
    // 状态显示
    statusLabel = new QLabel("状态: 未连接");
    statusLabel->setAlignment(Qt::AlignCenter);
    mainLayout->addWidget(statusLabel);
    
    // 容器列表
    containersList = new QListWidget();
    containersList->setStyleSheet(
        "QListWidget { background-color: rgba(240, 240, 240, 0.8); border-radius: 5px; }"
    );
    mainLayout->addWidget(containersList);
    
    // 按钮布局
    QHBoxLayout *buttonLayout = new QHBoxLayout();
    
    refreshButton = new QPushButton("刷新");
    refreshButton->setStyleSheet(
        "QPushButton { background-color: #4CAF50; color: white; border-radius: 5px; padding: 5px; }"
        "QPushButton:hover { background-color: #45a049; }"
    );
    buttonLayout->addWidget(refreshButton);
    
    hideButton = new QPushButton("隐藏");
    hideButton->setStyleSheet(
        "QPushButton { background-color: #f44336; color: white; border-radius: 5px; padding: 5px; }"
        "QPushButton:hover { background-color: #da190b; }"
    );
    buttonLayout->addWidget(hideButton);
    
    mainLayout->addLayout(buttonLayout);
    
    QVBoxLayout *outerLayout = new QVBoxLayout(this);
    outerLayout->addWidget(contentWidget);
    
    connect(hideButton, &QPushButton::clicked, this, &DesktopWidget::hide);
}

void DesktopWidget::updateStatus(const QString &status)
{
    statusLabel->setText("状态: " + status);
    
    if (status == "运行中") {
        statusLabel->setStyleSheet("color: green; font-weight: bold;");
    } else if (status == "停止") {
        statusLabel->setStyleSheet("color: red; font-weight: bold;");
    } else {
        statusLabel->setStyleSheet("color: black;");
    }
}

void DesktopWidget::updateContainers(const QList<QMap<QString, QString>> &containers)
{
    containersList->clear();
    
    for (const auto &container : containers) {
        QListWidgetItem *item = new QListWidgetItem();
        QString name = container.value("name", "未知");
        QString status = container.value("status", "未知");
        QString info = name + " - " + status;
        item->setText(info);
        
        if (status == "running") {
            item->setForeground(Qt::green);
        } else if (status == "stopped") {
            item->setForeground(Qt::red);
        }
        
        containersList->addItem(item);
    }
}

void DesktopWidget::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton) {
        lastMousePos = event->globalPos() - frameGeometry().topLeft();
        event->accept();
    }
}

void DesktopWidget::mouseMoveEvent(QMouseEvent *event)
{
    if (event->buttons() & Qt::LeftButton) {
        move(event->globalPos() - lastMousePos);
        event->accept();
    }
}
